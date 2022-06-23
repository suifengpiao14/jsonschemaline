package jsonschemaline

import (
	"strings"

	"github.com/go-errors/errors"
)

type DefaultJson struct {
	ID      string
	Version string
	Json    string
}

func ParseDefaultJson(linechemas string) (defaultJsons []*DefaultJson) {
	linechemas = strings.TrimSpace(strings.ReplaceAll(linechemas, "\r\n", EOF))
	arr := strings.Split(linechemas, EOF_DOUBLE)
	defaultJsons = make([]*DefaultJson, 0)
	for _, lineschema := range arr {
		defaultJson := ParseOneDefaultJson(lineschema)
		defaultJsons = append(defaultJsons, defaultJson)
	}
	return defaultJsons
}

func ParseOneDefaultJson(lineschema string) (defaultJson *DefaultJson) {
	defaultJson = new(DefaultJson)
	tagLineKVpairs := SplitMultilineSchema(lineschema)
	metaline, ok := GetMetaLine(tagLineKVpairs)
	if !ok {
		err := errors.Errorf("ParseOneDefaultJson meta line required,got: %#v", lineschema)
		panic(err)
	}
	meta := ParseMeta(*metaline)
	defaultJson.ID = meta.ID.String()
	defaultJson.Version = meta.Version
	fullnameList := make([]string, 0)
	for _, lineTags := range tagLineKVpairs {
		var (
			fullname     string
			src          string
			dst          string
			format       string
			defaultValue string
			required     bool
		)
		for _, kvPair := range lineTags {
			switch kvPair.Key {
			case "fullname":
				fullname = kvPair.Value
				fullnameList = append(fullnameList, fullname)
			case "src":
				src = kvPair.Value
			case "dst":
				dst = kvPair.Value
			case "format":
				format = kvPair.Value
			case "required":
				required = true
			case "default":
				defaultValue = kvPair.Value
			}
		}
		if fullname == "" {
			continue
		}
	}
	return defaultJson
}
