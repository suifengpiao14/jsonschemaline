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
		result := gjson.Get(data, kv.Key)
		if result.String() == "" || (result.Type == gjson.Number && result.Int() == 0) {
			data, err = sjson.Set(data, kv.Key, kv.Value)
			if err != nil {
				return "", err
			}
		}
	}
	return data, err
}
