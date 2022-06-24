package jsonschemaline

import (
	"fmt"
	"strings"

	"github.com/tidwall/sjson"
)

type DefaultJson struct {
	ID      string
	Version string
	Json    string
}

func ParseDefaultJson(linechemas string) (defaultJsons []*DefaultJson, err error) {
	linechemas = strings.TrimSpace(strings.ReplaceAll(linechemas, "\r\n", EOF))
	arr := strings.Split(linechemas, EOF_DOUBLE)
	defaultJsons = make([]*DefaultJson, 0)
	for _, lineschema := range arr {
		defaultJson, err := ParseOneDefaultJson(lineschema)
		if err != nil {
			return nil, err
		}
		defaultJsons = append(defaultJsons, defaultJson)
	}
	return defaultJsons, nil
}

func ParseOneDefaultJson(lineschema string) (defaultJson *DefaultJson, err error) {
	defaultJson = new(DefaultJson)
	jsonschemaline, err := ParseJsonschemaline(lineschema)
	if err != nil {
		return nil, err
	}
	id := jsonschemaline.Meta.ID.String()
	defaultJson.ID = id
	defaultJson.Version = jsonschemaline.Meta.Version
	kvmap := make(map[string]string)
	for _, item := range jsonschemaline.Items {
		if item.Default != "" {
			path := strings.ReplaceAll(item.Fullname, "[]", ".#")
			k := fmt.Sprintf("%s.%s", id, path)
			kvmap[k] = item.Default
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
