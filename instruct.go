package jsonschemaline

import (
	"fmt"
	"sort"
	"strings"
)

const (
	LINE_SCHEMA_DIRECTION_OUT     = "out"     // 出参
	LINE_SCHEMA_DIRECTION_IN      = "in"      // 入参
	LINE_SCHEMA_DIRECTION_CONVERT = "convert" // 内部转换
)

type Instruct struct {
	ID            string // 唯一标识
	Cmd           string
	Src           string
	Dst           string
	Tpl           string
	ExtraEndTpl   []string // 一个元素代表一个模板命令,方便去重
	ExtraStartTpl []string // 一个元素代表一个模板命令,方便去重
}

type Instructs []*Instruct

func (instructs Instructs) Unique() Instructs {
	out := Instructs{}
	m := make(map[string]bool)
	for _, instruct := range instructs {
		_, ok := m[instruct.ID]
		if !ok {
			m[instruct.ID] = true
			out = append(out, instruct)
		}
	}
	return out
}

func (a Instructs) Len() int      { return len(a) }
func (a Instructs) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Instructs) Less(i, j int) bool {
	return strings.Count(a[i].ID, ".") < strings.Count(a[j].ID, ".")
}

type InstructTpl struct {
	ID        string
	Version   string
	Type      string
	Instructs Instructs
}
type InstructTpls []*InstructTpl

func (instructTpls InstructTpls) GetByID(id string) *InstructTpl {
	for _, instruct := range instructTpls {
		if instruct.ID == id {
			out := instruct
			return out
		}
	}
	return nil
}

func (instructTpl *InstructTpl) String() string {
	allTplArr := make([]string, 0)
	allTplArr = append(allTplArr, fmt.Sprintf(`{{define "%s"}}`, instructTpl.ID))
	extraStartTpls := make(map[string]bool)
	middlTpls := make([]string, 0)
	extraEndTpls := make(map[string]bool)
	for _, instruct := range instructTpl.Instructs {

		for _, extraTpl := range instruct.ExtraStartTpl {
			extraStartTpls[extraTpl] = true
		}

		dst := instruct.Dst
		src := instruct.Src
		var value string
		if instruct.Tpl != "" {
			value = instruct.Tpl
		} else {
			value = fmt.Sprintf(`{{%s . "%s" "%s"}}`, instruct.Cmd, dst, src)
		}
		middlTpls = append(middlTpls, value)
		for _, extraTpl := range instruct.ExtraEndTpl {
			extraEndTpls[extraTpl] = true
		}
	}

	for extraTpl := range extraStartTpls {
		allTplArr = append(allTplArr, extraTpl)
	}
	allTplArr = append(allTplArr, middlTpls...)
	for extraTpl := range extraEndTpls {
		allTplArr = append(allTplArr, extraTpl)
	}

	allTplArr = append(allTplArr, `{{end}}`)
	out := strings.Join(allTplArr, "\n")
	return out
}

func ParseInstructTp(lineschema Jsonschemaline) (instructTpl *InstructTpl) {
	instructTpl = new(InstructTpl)
	instructTpl.ID = lineschema.Meta.ID
	instructTpl.Version = lineschema.Meta.Version
	instructTpl.Type = lineschema.Meta.Direction
	for _, item := range lineschema.Items {
		instruct := Instruct{
			ID:  item.Fullname,
			Src: item.Src,
			Dst: item.Dst,
		}

		if len(instruct.Src) > 2 && instruct.Src[:2] == "{{" {
			instruct.Tpl = instruct.Src
		}
		if len(instruct.Dst) > 2 && instruct.Dst[:2] == "{{" {
			instruct.Tpl = instruct.Dst
		}

		if instruct.Tpl == "" { // 本身不是tpl的情况下，构造tpl
			switch item.Format {
			case "number", "int", "integer", "float":
				instruct.Cmd = "getSetNumber"
			default:
				instruct.Cmd = "getSetValue"
			}
		}
		instructTpl.Instructs = append(instructTpl.Instructs, &instruct)
	}

	if instructTpl.Type == LINE_SCHEMA_DIRECTION_OUT {
		*instructTpl = FormatOutputTplInstruct(*instructTpl)
	}

	return instructTpl
}

func FormatOutputTplInstruct(instructTpl InstructTpl) (newInstructTpl InstructTpl) {
	newInstructTpl = instructTpl
	instructs := Instructs{}
	root := instructTpl.ID
	for _, instruct := range instructTpl.Instructs {
		newInstruct := FormatOutputInstruct(*instruct, root)
		instructs = append(instructs, &newInstruct)
	}
	instructs = instructs.Unique()
	sort.Sort(instructs)
	newInstructTpl.Instructs = instructs
	return newInstructTpl
}

func FormatOutputInstruct(instruct Instruct, root string) (newInstruct Instruct) {
	newInstruct = instruct
	fullname := instruct.ID
	for {
		if fullname == "" {
			break
		}
		lastIndex := strings.LastIndex(fullname, ".")
		if lastIndex < 0 {
			break
		}
		fullname = fullname[:lastIndex]
		if strings.HasSuffix(fullname, "[]") {
			fullname = strings.TrimSuffix(fullname, "[]")
			startTpl := fmt.Sprintf(`{{setValue . "%s.%s" list }}`, root, fullname) // 数组需要后续翻转，所以先声明为对象
			newInstruct.ExtraStartTpl = append(newInstruct.ExtraStartTpl, startTpl)
		} else {
			startTpl := fmt.Sprintf(`{{setValue . "%s.%s" dict }}`, root, fullname)
			newInstruct.ExtraStartTpl = append(newInstruct.ExtraStartTpl, startTpl)
		}
	}
	return newInstruct
}
