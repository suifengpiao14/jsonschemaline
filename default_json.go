package jsonschemaline

import (
	"strings"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
)

type DefaultJson struct {
	ID      string
	Version string
	Json    string
}

func ParseDefaultJson(jsonschemaline Jsonschemaline) (defaultJson *DefaultJson, err error) {
	defaultJson = new(DefaultJson)
	id := jsonschemaline.Meta.ID.String()
	defaultJson.ID = id
	defaultJson.Version = jsonschemaline.Meta.Version
	kvmap := make(map[string]string)
	for _, item := range jsonschemaline.Items {
		if item.Default != "" || item.AllowEmptyValue {
			path := strings.ReplaceAll(item.Fullname, "[]", ".#")
			//k := fmt.Sprintf("%s.%s", id, path)
			kvmap[path] = item.Default
		}
	}
	jsonContent := ""
	for k, v := range kvmap {
		jsonContent, err = sjson.Set(jsonContent, k, v)
		if err != nil {
			return nil, err
		}
	}
	defaultJson.Json = jsonContent
	return defaultJson, nil
}

func JsonMerge(first string, second string, more ...string) (merge string, err error) {
	combinedPatch := []byte(second)
	for _, patch := range more {
		combinedPatch, err = jsonpatch.MergeMergePatches(combinedPatch, []byte(patch))
		if err != nil {
			err = errors.WithStack(err)
			return "", err
		}

	}
	mb, err := jsonpatch.MergePatch([]byte(first), combinedPatch)
	if err != nil {
		err = errors.WithStack(err)
		return "", err
	}
	merge = string(mb)
	return merge, err
}
