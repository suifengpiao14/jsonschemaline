package jsonschemaline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"github.com/suifengpiao14/kvstruct"
)

func JsonSchema2LineSchema(jsonschema string) (lineschema *Jsonschemaline, err error) {
	var schema map[string]interface{}
	err = json.Unmarshal([]byte(jsonschema), &schema)
	if err != nil {
		return nil, err
	}
	kvs := jsonSchema2KVS(schema, "")
	version, _ := kvs.GetFirstByKey("$schema")
	id, _ := kvs.GetFirstByKey("$id")
	if id.Value == "" {
		id.Value = "example"
	}
	kvs1 := dealPropertiesAndItemsdealRequired(kvs)
	kvs2 := dealRequired(kvs1)
	m := make(map[string][][2]string)
	for _, kv := range kvs2 {
		if strings.HasPrefix(kv.Key, "$") {
			continue
		}
		lastDot := strings.LastIndex(strings.Trim(kv.Key, "."), ".")
		fullname := ""
		key := kv.Key
		if lastDot > -1 {
			fullname, key = kv.Key[:lastDot], kv.Key[lastDot+1:]
		}
		if _, ok := m[fullname]; !ok {
			m[fullname] = make([][2]string, 0)
		}

		m[fullname] = append(m[fullname], [2]string{key, kv.Value})
	}

	var w bytes.Buffer
	w.WriteString(fmt.Sprintf("version=%s,direction=%s,id=%s\n", version.Value, "in", id.Value))
	for fullname, linePairs := range m {
		if fullname == "" {
			continue
		}
		pairs := make([]string, 0)
		pairs = append(pairs, fmt.Sprintf("fullname=%s", fullname), fmt.Sprintf("dst=%s", fullname))
		for _, pair := range linePairs {
			pairs = append(pairs, strings.Join(pair[:], "="))
		}
		w.WriteString(strings.Join(pairs, ","))
		w.WriteString("\n")
	}

	lineschemastr := w.String()

	lineschema, err = ParseJsonschemaline(lineschemastr)
	if err != nil {
		return nil, err
	}

	return lineschema, nil
}

// jsonschema 转 kvs
func jsonSchema2KVS(schema map[string]interface{}, prefix string) kvstruct.KVS {
	kvs := kvstruct.KVS{}
	for key, value := range schema {
		fieldName := fmt.Sprintf("%s%s", prefix, key)
		switch valueType := value.(type) {
		case map[string]interface{}:
			// 递归处理子对象
			kvs.Add(jsonSchema2KVS(valueType, fieldName+".")...)
		case string:
			kvs.Add(kvstruct.KV{
				Key:   fieldName,
				Value: valueType,
			})
		case []interface{}:
			b, _ := json.Marshal(value)
			valueStr := string(b)
			kvs.Add(kvstruct.KV{
				Key:   fieldName,
				Value: valueStr,
			})
		default:
			kvs.Add(kvstruct.KV{
				Key:   fieldName,
				Value: cast.ToString(value),
			})
		}
	}
	return kvs
}

func dealPropertiesAndItemsdealRequired(kvs kvstruct.KVS) (newKvs kvstruct.KVS) {
	newKvs = kvstruct.KVS{}
	keywordItem := "items."
	tmpKvs := make(kvstruct.KVS, 0)
	for _, kv := range kvs {
		segments := strings.Split(kv.Key, keywordItem)
		prefix := ""
		for i, segment := range segments {
			parent := fmt.Sprintf("%s%s", prefix, segment)
			parentType := fmt.Sprintf("%stype", parent)
			parentTypeKv, _ := kvs.GetFirstByKey(parentType)
			if parentTypeKv.Value == "array" {
				segments[i] = "[]"
			}
		}
		key := strings.Join(segments, "")
		tmpKvs.Add(kvstruct.KV{Key: key, Value: kv.Value})
	}

	keywordProperties := "properties."
	for _, kv := range kvs {
		segments := strings.Split(kv.Key, keywordProperties)
		prefix := ""
		for i, segment := range segments {
			parent := fmt.Sprintf("%s%s", prefix, segment)
			parentType := fmt.Sprintf("%stype", parent)
			parentTypeKv, _ := kvs.GetFirstByKey(parentType)
			if parentTypeKv.Value == "object" {
				segments[i] = ""
			}
		}
		key := strings.Join(segments, "")
		newKvs.Add(kvstruct.KV{Key: key, Value: kv.Value})
	}
	return newKvs
}

func dealRequired(kvs kvstruct.KVS) (newKvs kvstruct.KVS) {
	newKvs = make(kvstruct.KVS, 0)

	for _, kv := range kvs {
		requiredLastIndex := strings.LastIndex(kv.Key, "required")
		if requiredLastIndex < 0 {
			newKvs.Add(kv)
			continue
		}
		prefix := kv.Key[:requiredLastIndex]
		typeKv, _ := kvs.GetFirstByKey(fmt.Sprintf("%stype", prefix))
		if typeKv.Value != "array" && typeKv.Value != "object" {
			newKvs.Add(kv)
			continue
		}
		var keys = make([]string, 0)
		err := json.Unmarshal([]byte(kv.Value), &keys)
		if err != nil {
			panic(err)
		}
		for _, k := range keys {
			newKvs.Add(kvstruct.KV{
				Key:   fmt.Sprintf("%s%s.required", prefix, k),
				Value: "true",
			})
		}
	}
	return newKvs
}
