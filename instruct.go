package jsonschemaline

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
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
	out := new(InstructTpl)
	for _, instruct := range instructTpls {
		if instruct.ID == ID(id) {
			out = instruct
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
		if instructTpl.Type == INSTRUCT_TYPE_OUT { // 输出时，dst 需要加上 根 元素
			dst = fmt.Sprintf("%s.%s", instructTpl.ID, dst)
		}
		if instructTpl.Type == INSTRUCT_TYPE_IN { // 输入时, src 需要加上 跟 元素
			src = fmt.Sprintf("%s.%s", instructTpl.ID, src)
		}
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

func ParseInstructTp(lineschemas string) (instructTpls InstructTpls) {

	lineschemas = strings.TrimSpace(strings.ReplaceAll(lineschemas, "\r\n", EOF))
	arr := strings.Split(lineschemas, EOF_DOUBLE)
	instructTpls = InstructTpls{}
	for _, lineschema := range arr {
		instructTpl := ParseOneInstructTp(lineschema)
		instructTpls = append(instructTpls, instructTpl)
	}
	return instructTpls
}

func ParseOneInstructTp(lineschema string) (instructTpl *InstructTpl) {
	instructTpl = new(InstructTpl)
	tagLineKVpairs := SplitMultilineSchema(lineschema)
	metaline, ok := GetMetaLine(tagLineKVpairs)
	if !ok {
		err := errors.Errorf("meta line required,got: %#v", lineschema)
		panic(err)
	}
	meta := ParseMeta(*metaline)
	instructTpl.ID = meta.ID
	instructTpl.Version = meta.Version
	fullnameList := make([]string, 0)
	for _, lineTags := range tagLineKVpairs {
		var (
			fullname string
			src      string
			dst      string
			format   string
			instruct Instruct
		)
		for _, kvPair := range lineTags {
			switch kvPair.Key {
			case "fullname":
				fullname = kvPair.Value
				fullnameList = append(fullnameList, fullname)
			case "src":
				src = kvPair.Value
			case "dst":
				dst = kvPair.Value
			case "format":
				format = kvPair.Value
			}
		}
		if fullname == "" {
			continue
		}
		srcOrDst := strings.ReplaceAll(fullname, "[]", ".#")
		if src == "" {
			instructTpl.Type = INSTRUCT_TYPE_IN
			src = srcOrDst
		}
		if dst == "" {
			instructTpl.Type = INSTRUCT_TYPE_OUT // dst 为空，则说明fullname 充当目标地址，说明src不为空，即复制到json
			dst = srcOrDst
		}
		if src == dst {
			continue
		}
		instruct.Src = src
		instruct.Dst = dst
		if len(instruct.Src) > 2 && instruct.Src[:2] == "{{" {
			instruct.Tpl = instruct.Src
		}
		if len(instruct.Dst) > 2 && instruct.Dst[:2] == "{{" {
			instruct.Tpl = instruct.Dst
		}

		if instruct.Tpl == "" { // 本身不是tpl的情况下，构造tpl
			switch format {
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
