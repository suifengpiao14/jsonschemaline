package jsonschemaline

import (
	"fmt"
	"sort"
	"strings"
)

const (
	INSTRUCT_TYPE_OUT     = "out"     // 入参
	INSTRUCT_TYPE_IN      = "in"      // 出参
	INSTRUCT_TYPE_CONVERT = "convert" // 内部转换
)

type Instruct struct {
	ID  string // 唯一标识
	Cmd string
	Src string
	Dst string
	Tpl string
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
	ID        ID
	Version   string
	Type      string
	Instructs Instructs
}
type InstructTpls []*InstructTpl

func (instructTpls InstructTpls) GetByID(id string) *InstructTpl {
	for _, instruct := range instructTpls {
		if instruct.ID == ID(id) {
			out := instruct
			return out
		}
	}
	return nil
}

func (instructTpl *InstructTpl) String() string {
	valueArr := make([]string, 0)
	valueArr = append(valueArr, fmt.Sprintf(`{{define "%s"}}`, instructTpl.ID))
	for _, instruct := range instructTpl.Instructs {
		dst := instruct.Dst
		src := instruct.Src
		var value string
		if instruct.Tpl != "" {
			value = instruct.Tpl
		} else {
			value = fmt.Sprintf(`{{%s . "%s" "%s"}}`, instruct.Cmd, dst, src)
		}
		valueArr = append(valueArr, value)
	}
	valueArr = append(valueArr, `{{end}}`)
	out := strings.Join(valueArr, "\n")
	return out
}

func ParseInstructTp(lineschema Jsonschemaline) (instructTpl *InstructTpl) {
	instructTpl = new(InstructTpl)
	instructTpl.ID = lineschema.Meta.ID
	instructTpl.Version = lineschema.Meta.Version
	instructTpl.Type = lineschema.Meta.Direction
	fullnameList := make([]string, 0)
	for _, item := range lineschema.Items {
		instruct := Instruct{
			ID:  item.Fullname,
			Src: item.Src,
			Dst: item.Dst,
		}
		fullnameList = append(fullnameList, item.Fullname)

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

	if instructTpl.Type == INSTRUCT_TYPE_OUT {
		parentInstructs := parseParentInstruct(fullnameList, instructTpl.ID.String())
		parentInstructs = append(parentInstructs, instructTpl.Instructs...) // 确保跟元素先初始化（防止空数据时，根元素不输出）
		instructTpl.Instructs = parentInstructs
	}

	return instructTpl
}

func parseParentInstruct(fullnames []string, root string) (instructs Instructs) {
	instructs = Instructs{}
	for _, fullname := range fullnames {
		instrList := parseFullname(fullname, root)
		instructs = append(instructs, instrList...)
	}
	instructs = instructs.Unique()
	sort.Sort(instructs)
	return instructs
}

func parseFullname(fullname string, root string) (instructs Instructs) {
	instructs = Instructs{}
	for {
		if fullname == "" {
			break
		}
		lastIndex := strings.LastIndex(fullname, ".")
		if lastIndex < 0 {
			break
		}
		fullname = fullname[:lastIndex]

		instruct := &Instruct{ID: fullname}
		if strings.HasSuffix(fullname, "[]") {
			fullname = strings.TrimSuffix(fullname, "[]")
			instruct.ID = fullname
			instruct.Tpl = fmt.Sprintf(`{{setValue . "%s.%s" list }}`, root, fullname)
		} else {
			instruct.Tpl = fmt.Sprintf(`{{setValue . "%s.%s" dict }}`, root, fullname)
		}
		instructs = append(instructs, instruct)
	}
	return instructs
}
