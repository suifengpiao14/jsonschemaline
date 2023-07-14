package jsonschemaline

import (
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/suifengpiao14/kvstruct"
)

func JsonMergeIgnoreEmptyString(first string, second string) (merge string, err error) {
	patch, err := jsonpatch.CreateMergePatch([]byte(first), []byte(second))
	if err != nil {
		return "", err
	}
	kvs := kvstruct.JsonToKVS(string(patch), "")
	fielterKvs := make(kvstruct.KVS, 0)
	for _, kv := range kvs {
		if kv.Value == "" {
			continue
		}
		fielterKvs.Add(kv)
	}
	patchStr, err := fielterKvs.Json(false)
	if err != nil {
		return "", err
	}

	mergeb, err := jsonpatch.MergePatch([]byte(first), []byte(patchStr))
	if err != nil {
		return "", err
	}
	merge = string(mergeb)
	return merge, err
}
