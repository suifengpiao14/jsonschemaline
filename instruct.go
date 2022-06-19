package jsonschemaline

import (
	"fmt"
	"strings"
)

const (
	INSTRUCT_COPY_2_JSON    = "copy2json"
	INSTRUCT_COPY_2_CONTEXT = "copy2context"
)

type Instruct struct {
	Name string
	Cmd  string
	Src  string
	Dst  string
}

type Instructs []*Instruct

func (instructs Instructs) GetByName(name string) Instructs {
	out := Instructs{}
	for _, instruct := range instructs {
		if instruct.Name == name {
			out = append(out, instruct)
		}
	}
	return out
}

func (instructs Instructs) String(root string) string {
	valueArr := make([]string, 0)
	for _, instruct := range instructs {
		dst := instruct.Dst
		src := instruct.Src
		if instruct.Name == INSTRUCT_COPY_2_CONTEXT {
			src = fmt.Sprintf("%s.%s", root, src)
		}
		if instruct.Name == INSTRUCT_COPY_2_JSON {
			dst = fmt.Sprintf("%s.%s", root, dst)
		}
		value := fmt.Sprintf(`{{%s . "%s" %s""}}`, instruct.Cmd, dst, src)
		valueArr = append(valueArr, value)
	}
	out := strings.Join(valueArr, "\n")
	return out
}

func ParseInstruct(lineschema string) (instructArr Instructs) {
	multilineTags := SplitMultilineSchema(lineschema)
	instructArr = make(Instructs, 0)
	for _, lineTags := range multilineTags {
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
			instruct.Name = INSTRUCT_COPY_2_JSON
			src = srcOrDst
		}
		if dst == "" {
			instruct.Name = INSTRUCT_COPY_2_CONTEXT
			dst = srcOrDst
		}
		if src == dst {
			continue
		}
		switch format {
		case "number", "int", "integer", "float":
			instruct.Cmd = "getSetNumber"
		default:
			instruct.Cmd = "getSetValue"
		}
		instruct.Src = src
		instruct.Dst = dst
		instructArr = append(instructArr, &instruct)
	}
	return instructArr
}
