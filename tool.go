package jsonschemaline

import (
	"github.com/suifengpiao14/kvstruct"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func MergeDefault(data string, defaul string) (merge string, err error) {
	kvs := kvstruct.JsonToKVS(defaul, "")
	for _, kv := range kvs {
		if kv.Value == "" {
			continue
		}
		v := gjson.Get(data, kv.Key).String()
		if v == "" {
			data, err = sjson.Set(data, kv.Key, kv.Value)
			if err != nil {
				return "", err
			}
		}
	}
	return data, err
}
