package jsonschemaline

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/suifengpiao14/kvstruct"
)

func JsonSchema2LineSchema(jsonschema string) (lineschema *Jsonschemaline, err error) {
	kvs := kvstruct.JsonToKVS(jsonschema, "")
	version, _ := kvs.GetFirstByKey("$schema")
	id, _ := kvs.GetFirstByKey("$id")
	lineschema.Meta = &Meta{
		Version:   version.Value,
		ID:        ID(id.Value),
		Direction: "in",
	}
	kvs1 := dealPropertiesAndItemsdealRequired(kvs)
	kvs2 := dealRequired(kvs1)
	fmt.Println(kvs2)
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

	return
}
func dealPropertiesAndItemsdealRequired(kvs kvstruct.KVS) (newKvs kvstruct.KVS) {
	replacer := strings.NewReplacer("properties.", "", "items.", "[]")
	newKvs = make(kvstruct.KVS, 0)
	for _, kv := range kvs {
		newKv := kvstruct.KV{
			Key:   replacer.Replace(kv.Key),
			Value: kv.Value,
		}
		newKvs.Add(newKv)
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
