package jsonschemaline

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cast"
	"github.com/suifengpiao14/kvstruct"
)

// jsonschema 转 kvs
func jsonSchema2KVS(schema map[string]interface{}, prefix string) kvstruct.KVS {
	kvs := kvstruct.KVS{}
	for key, value := range schema {
		fieldName := fmt.Sprintf("%s%s", prefix, key)
		switch valueType := value.(type) {
		case map[string]interface{}:
			// 递归处理子对象
			kvs.Add(jsonSchema2KVS(valueType, fieldName+".")...)
		default:
			// 将键值对格式化为字符串
			kvs.Add(kvstruct.KV{
				Key:   fieldName,
				Value: cast.ToString(value),
			})
		}
	}
	return kvs
}

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
	lineschema.Meta = &Meta{
		Version:   version.Value,
		ID:        ID(id.Value),
		Direction: "in",
	}
	kvs1 := dealPropertiesAndItemsdealRequired(kvs)
	kvs2 := dealRequired(kvs1)
	kvs3 := dealArray(kvs2)
	fmt.Println(kvs3)
	return
}
func dealArray(kvs kvstruct.KVS) (newKvs kvstruct.KVS) {
	newKvs = make(kvstruct.KVS, 0)
	exp := regexp.MustCompile(`^(.*)\.\d+$`)
	for _, kv := range kvs {
		match := exp.FindStringSubmatch(kv.Key)
		if len(match) == 0 { // 非数组列复制,不修改
			newKvs.Add(kv)
		}
		key := match[1]
		if key == "" {
			continue
		}
		newkv, _ := newKvs.GetFirstByKey(key)
		var value = kv.Value
		if newkv.Value != "" {
			value = fmt.Sprintf("%s.%s", newkv.Value, value)

		}
		newkv.Key = key
		newkv.Value = value
		newKvs.AddReplace(newkv)
	}

	return newKvs
}
func dealPropertiesAndItemsdealRequired(kvs kvstruct.KVS) (newKvs kvstruct.KVS) {
	arrayKeyReplaceKvs := kvstruct.KVS{}
	objectKeys := make([]string, 0)
	typeExt := ".type"

	newKvs = make(kvstruct.KVS, 0)
	for _, kv := range kvs {
		if strings.ToLower(kv.Value) == "array" && strings.HasSuffix(kv.Key, typeExt) {

			arrayKeys = append(arrayKeys, fmt.Sprintf("%s.items", strings.TrimSuffix(kv.Key, typeExt)))
		}
	}
	for _, kv := range kvs {

	}

	return newKvs
}

func dealRequired(kvs kvstruct.KVS) (newKvs kvstruct.KVS) {
	newKvs = make(kvstruct.KVS, 0)
	exp := regexp.MustCompile(`^(.*)required\.\d+$`)
	for _, kv := range kvs {
		match := exp.FindStringSubmatch(kv.Key)
		if len(match) == 0 { // 非required列复制,不修改
			newKvs.Add(kv)
		}
		prefix := match[1]
		if prefix != "" {
			prefix = fmt.Sprintf("%s.", prefix)
		}
		newKvs.Add(kvstruct.KV{Key: fmt.Sprintf("%s%s", prefix, kv.Value), Value: "true"})
	}
	return newKvs
}

/**
[{ \$schema http://json-schema.org/draft-07/schema#}
{ \$id execAPIInquiryScreenIdentifyUpdate}
{ type object}
{ required.0 config}
{ config.id.type string}
{ config.id.format number}
{ config.status.type string}
{ config.status.enum.0 1}
{ config.status.enum.1 2}
{ config.identify.type string}
{ config.identify.minLength 1}
{ config.merchantId.type string}
{ config.merchantId.format number}
{ config.merchantName.type string}
{ config.merchantName.minLength 1}
{ config.operateName.type string}
{ config.operateName.minLength 1}
{ config.storeId.type string}
{ config.storeId.format number}
{ config.storeName.type string}
{ config.storeName.minLength 1}
{ config.type object}
{ config.required.0 id}
{ config.required.1 status}
{ config.required.2 identify}
{ config.required.3 merchantId}
{ config.required.4 merchantName}
{ config.required.5 operateName}
{ config.required.6 storeId}
{ config.required.7 storeName}]
**/
