package jsonschemaline

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-errors/errors"
)

type JsonschemalineItem struct {
	Comments string `json:"comment,omitempty"` // section 8.3

	Type             string        `json:"type,omitempty"`                    // section 6.1.1
	Enum             []interface{} `json:"enum,omitempty"`                    // section 6.1.2
	Const            interface{}   `json:"const,omitempty"`                   // section 6.1.3
	MultipleOf       int           `json:"multipleOf,omitempty,string"`       // section 6.2.1
	Maximum          int           `json:"maximum,omitempty,string"`          // section 6.2.2
	ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty,string"` // section 6.2.3
	Minimum          int           `json:"minimum,omitempty,string"`          // section 6.2.4
	ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty,string"` // section 6.2.5
	MaxLength        int           `json:"maxLength,omitempty,string"`        // section 6.3.1
	MinLength        int           `json:"minLength,omitempty,string"`        // section 6.3.2
	Pattern          string        `json:"pattern,omitempty"`                 // section 6.3.3
	MaxItems         int           `json:"maxItems,omitempty,string"`         // section 6.4.1
	MinItems         int           `json:"minItems,omitempty,string"`         // section 6.4.2
	UniqueItems      bool          `json:"uniqueItems,omitempty,string"`      // section 6.4.3
	MaxContains      uint          `json:"maxContains,omitempty,string"`      // section 6.4.4
	MinContains      uint          `json:"minContains,omitempty,string"`      // section 6.4.5
	MaxProperties    int           `json:"maxProperties,omitempty,string"`    // section 6.5.1
	MinProperties    int           `json:"minProperties,omitempty,string"`    // section 6.5.2
	Required         bool          `json:"required,omitempty,string"`         // section 6.5.3
	// RFC draft-bhutton-json-schema-validation-00, section 7
	Format string `json:"format,omitempty"`
	// RFC draft-bhutton-json-schema-validation-00, section 8
	ContentEncoding  string `json:"contentEncoding,omitempty"`   // section 8.3
	ContentMediaType string `json:"contentMediaType,omitempty"`  // section 8.4
	Title            string `json:"title,omitempty"`             // section 9.1
	Description      string `json:"description,omitempty"`       // section 9.1
	Default          string `json:"default,omitempty"`           // section 9.2
	Deprecated       bool   `json:"deprecated,omitempty,string"` // section 9.3
	ReadOnly         bool   `json:"readOnly,omitempty,string"`   // section 9.4
	WriteOnly        bool   `json:"writeOnly,omitempty,string"`  // section 9.4
	Example          string `json:"examples,omitempty,string"`   // section 9.5
	Src              string `json:"src,omitempty"`
	Dst              string `json:"dst,omitempty"`
	Fullname         string `json:"fullname"`
}
type Meta struct {
	ID      ID     `json:"id"`
	Version string `json:"version"`
}

func IsMetaLine(lineTags TagLineKVpair) bool {
	hasFullname, hasId := false, false
	for _, kvPair := range lineTags {
		switch kvPair.Key {
		case "id":
			hasId = true
		case "fullname":
			hasFullname = true
		}
	}
	is := hasId && !hasFullname
	return is
}

type Jsonschemaline struct {
	Meta  *Meta
	Items []*JsonschemalineItem
}

type KVpair struct {
	Key   string
	Value string
}

type TagLineKVpair []KVpair

//ParseMultiJsonSchemaline 解析多个 jsonschemaline
func ParseMultiJsonSchemaline(jsonschemalineBlocks string) (jsonschemalines []*Jsonschemaline, err error) {
	jsonschemalineBlocks = strings.TrimSpace(strings.ReplaceAll(jsonschemalineBlocks, "\r\n", EOF))
	arr := strings.Split(jsonschemalineBlocks, EOF_DOUBLE)
	jsonschemalines = make([]*Jsonschemaline, 0)
	for _, lineschemaBlock := range arr {
		jsonschemaline, err := ParseJsonschemaline(lineschemaBlock)
		if err != nil {
			return nil, err
		}
		jsonschemalines = append(jsonschemalines, jsonschemaline)
	}
	return jsonschemalines, nil
}

//ParseJsonschemaline 解析单个 jsonschemaline
func ParseJsonschemaline(jsonschemalineBlock string) (jsonschemaline *Jsonschemaline, err error) {
	jsonschemaline = new(Jsonschemaline)
	jsonschemalineBlock = strings.TrimSpace(strings.ReplaceAll(jsonschemalineBlock, "\r\n", EOF))
	arr := strings.Split(jsonschemalineBlock, EOF)
	for _, raw := range arr {
		raw = strings.TrimSpace(raw)
		meta, item, err := ParseJsonschemalineRaw(raw)
		if err != nil {
			return nil, err
		}
		if meta != nil {
			jsonschemaline.Meta = meta
		} else if item != nil {
			jsonschemaline.Items = append(jsonschemaline.Items, item)
		}
	}
	if jsonschemaline.Meta == nil {
		err := errors.Errorf("jsonschemaline ID required,got:%s", jsonschemalineBlock)
		return nil, err
	}
	return jsonschemaline, nil
}

// ParseJsonschemalineRaw 解析 jsonschemaline 一行数据
func ParseJsonschemalineRaw(jsonschemalineRaw string) (meta *Meta, jsonschemalineItem *JsonschemalineItem, err error) {
	jsonschemalineRaw = PretreatJsonschemalineRaw(jsonschemalineRaw)
	kvStrArr := SplitOnUnescapedCommas(jsonschemalineRaw)
	kvMap := make(map[string]string)
	for _, kvStr := range kvStrArr {
		kvPair := strings.SplitN(kvStr, "=", 2)
		if len(kvPair) == 2 {
			k, v := strings.TrimSpace(kvPair[0]), strings.TrimSpace(kvPair[1])
			kvMap[k] = v
		}
	}
	_, hasId := kvMap["id"]
	_, hasFullname := kvMap["fullname"]
	jb, err := json.Marshal(kvMap)
	if err != nil {
		return nil, nil, err
	}
	if hasId && !hasFullname {
		meta = new(Meta)
		err = json.Unmarshal(jb, meta)
		if err != nil {
			return nil, nil, err
		}
		return meta, nil, nil
	}
	jsonschemalineItem = new(JsonschemalineItem)
	err = json.Unmarshal(jb, jsonschemalineItem)
	if err != nil {
		return nil, nil, err
	}

	return nil, jsonschemalineItem, nil
}

//PretreatJsonschemalineRaw 处理enum []格式
func PretreatJsonschemalineRaw(tag string) (formatTag string) {
	preg := "enum=\\[(.*)\\],"
	formatTag = strings.Trim(tag, ",")
	reg := regexp.MustCompile(preg)
	matchArr := reg.FindAllStringSubmatch(tag, -1)
	if len(matchArr) > 0 {
		replaceStr := "enum="
		for _, matchRaw := range matchArr {
			raw := strings.ReplaceAll(matchRaw[1], `"`, "")
			valArr := strings.Split(raw, ",")
			replaceStr = fmt.Sprintf("enum=%s,", strings.Join(valArr, ",enum="))
			formatTag = strings.ReplaceAll(formatTag, matchRaw[0], replaceStr)
		}
	}

	hasType := false
	kvStrArr := strings.Split(formatTag, ",")
	tmpArr := make([]string, 0)
	for _, kvStr := range kvStrArr {
		kvStr = strings.TrimSpace(kvStr)
		kvPair := strings.SplitN(kvStr, "=", 2)
		if len(kvPair) == 1 {
			kvPair = append(kvPair, "true") // bool 类型为true时，可以简写只有key
		}
		k, v := strings.TrimSpace(kvPair[0]), strings.TrimSpace(kvPair[1])
		hasType = hasType || k == "type"
		tmpArr = append(tmpArr, fmt.Sprintf("%s=%s", k, v))
	}
	if !hasType {
		tmpArr = append(tmpArr, "type=string") // 增加默认type=string
	}
	formatTag = strings.Join(tmpArr, ",")
	return formatTag
}
