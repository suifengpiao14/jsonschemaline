package jsonschemaline

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/suifengpiao14/kvstruct"
)

const (
	TOKEN_BEGIN = ','
	TOKEN_END   = '='
)

type KeyIndex struct {
	BeginAt int
	EndAt   int
}

// ParserLine 解析一行数据
func ParserLine(line string) (kvs kvstruct.KVS) {
	if line == "" {
		return nil
	}
	replacer := strings.NewReplacer(" ", "", "\n", "", "\t", "", "\r", "")
	line = replacer.Replace(line)
	ret := make([]string, 0)
	separated := strings.Split(line, ",")
	ret = append(ret, separated[0])
	i := 0
	for _, nextTag := range separated[1:] {
		if isToken(nextTag) {
			ret = append(ret, nextTag)
			i++
		} else {
			ret[i] = fmt.Sprintf("%s,%s", ret[i], nextTag)
		}
	}
	kvs = make(kvstruct.KVS, 0)
	for _, pair := range ret {
		arr := strings.SplitN(pair, "=", 2)
		if len(arr) == 1 {
			arr = append(arr, "true")
		}
		k, v := arr[0], arr[1]
		kv := kvstruct.KV{
			Key:   k,
			Value: v,
		}
		kvs.Add(kv)

	}

	fmt.Println(kvs)

	return
}
func isToken(s string) (yes bool) {
	for _, token := range getTokens() {
		yes = strings.HasPrefix(s, token)
		if yes {
			return yes
		}
	}
	return false
}

func getTokens() (tokens []string) {
	tokens = make([]string, 0)
	meta := new(Meta)
	var rt reflect.Type
	rt = reflect.TypeOf(meta).Elem()
	tokens = append(tokens, getJsonTagname(rt)...)
	item := new(JsonschemalineItem)
	rt = reflect.TypeOf(item).Elem()
	tokens = append(tokens, getJsonTagname(rt)...)

	return tokens
}

func getJsonTagname(rt reflect.Type) (jsonNames []string) {
	jsonNames = make([]string, 0)
	for i := 0; i < rt.NumField(); i++ {
		jsonTag := rt.Field(i).Tag.Get("json")
		index := strings.Index(jsonTag, ",")
		if index > 0 {
			jsonTag = jsonTag[:index]
		}
		jsonTag = strings.TrimSpace(jsonTag)
		if jsonTag != "-" {
			jsonNames = append(jsonNames, jsonTag)
		}
	}
	return jsonNames
}
