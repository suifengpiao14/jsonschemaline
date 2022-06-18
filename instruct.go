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
	Name  string
	Value string
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

func (instructs Instructs) String() string {
	valueArr := make([]string, 0)
	for _, instruct := range instructs {
		valueArr = append(valueArr, instruct.Value)
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
		var value string
		switch format {
		case "number", "int", "integer", "float":
			value = fmt.Sprintf(`{{getSetNumber . "%s" "%s"}}`, dst, src)
		default:
			value = fmt.Sprintf(`{{getSetValue . "%s" "%s"}}`, dst, src)
		}
		instruct.Value = value
		instructArr = append(instructArr, &instruct)
	}
	return instructArr
}
