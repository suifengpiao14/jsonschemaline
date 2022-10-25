package jsonschemaline

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
)

type JsonschemalineItem struct {
	Comments string `json:"comment,omitempty"` // section 8.3

	Type             string `json:"type,omitempty"`                    // section 6.1.1
	Enum             string `json:"enum,omitempty"`                    // section 6.1.2
	Const            string `json:"const,omitempty"`                   // section 6.1.3
	MultipleOf       int    `json:"multipleOf,omitempty,string"`       // section 6.2.1
	Maximum          int    `json:"maximum,omitempty,string"`          // section 6.2.2
	ExclusiveMaximum bool   `json:"exclusiveMaximum,omitempty,string"` // section 6.2.3
	Minimum          int    `json:"minimum,omitempty,string"`          // section 6.2.4
	ExclusiveMinimum bool   `json:"exclusiveMinimum,omitempty,string"` // section 6.2.5
	MaxLength        int    `json:"maxLength,omitempty,string"`        // section 6.3.1
	MinLength        int    `json:"minLength,omitempty,string"`        // section 6.3.2
	Pattern          string `json:"pattern,omitempty"`                 // section 6.3.3
	MaxItems         int    `json:"maxItems,omitempty,string"`         // section 6.4.1
	MinItems         int    `json:"minItems,omitempty,string"`         // section 6.4.2
	UniqueItems      bool   `json:"uniqueItems,omitempty,string"`      // section 6.4.3
	MaxContains      uint   `json:"maxContains,omitempty,string"`      // section 6.4.4
	MinContains      uint   `json:"minContains,omitempty,string"`      // section 6.4.5
	MaxProperties    int    `json:"maxProperties,omitempty,string"`    // section 6.5.1
	MinProperties    int    `json:"minProperties,omitempty,string"`    // section 6.5.2
	Required         bool   `json:"required,omitempty,string"`         // section 6.5.3
	// RFC draft-bhutton-json-schema-validation-00, section 7
	Format string `json:"format,omitempty"`
	// RFC draft-bhutton-json-schema-validation-00, section 8
	ContentEncoding  string        `json:"contentEncoding,omitempty"`   // section 8.3
	ContentMediaType string        `json:"contentMediaType,omitempty"`  // section 8.4
	Title            string        `json:"title,omitempty"`             // section 9.1
	Description      string        `json:"description,omitempty"`       // section 9.1
	Default          string        `json:"default,omitempty"`           // section 9.2
	Deprecated       bool          `json:"deprecated,omitempty,string"` // section 9.3
	ReadOnly         bool          `json:"readOnly,omitempty,string"`   // section 9.4
	WriteOnly        bool          `json:"writeOnly,omitempty,string"`  // section 9.4
	Example          string        `json:"example,omitempty,string"`    // section 9.5
	Src              string        `json:"src,omitempty"`
	Dst              string        `json:"dst,omitempty"`
	Fullname         string        `json:"fullname"`
	TagLineKVpair    TagLineKVpair `json:"-"`
}

var jsonschemalineItemOrder = []string{
	"fullname", "src", "dst", "type", "format", "pattern", "enum", "required", "title", "description", "default", "comment", "example", "deprecated", "const",
	"multipleOf", "maximum", "exclusiveMaximum", "minimum", "exclusiveMinimum", "maxLength", "minLength",
	"maxItems",
	"minItems",
	"uniqueItems",
	"maxContains",
	"minContains",
	"maxProperties",
	"minProperties",
	"contentEncoding",
	"contentMediaType",
	"readOnly",
	"writeOnly",
}

type Meta struct {
	ID        ID     `json:"id"`
	Version   string `json:"version"`
	Direction string `json:"direction"`
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

func (l *Jsonschemaline) String() string {

	lineArr := make([]string, 0)
	lineArr = append(lineArr, fmt.Sprintf("version=%s,id=%s", l.Meta.Version, l.Meta.ID))
	var linemap []map[string]string
	b, err := json.Marshal(l.Items)
	if err != nil {
		err = errors.WithStack(err)
		panic(err)
	}
	err = json.Unmarshal(b, &linemap)
	if err != nil {
		err = errors.WithStack(err)
		panic(err)
	}

	for _, m := range linemap {
		kvArr := make([]string, 0)
		for _, k := range jsonschemalineItemOrder {
			if l.Meta.Direction == LINE_SCHEMA_DIRECTION_IN && k == "src" {
				continue
			}
			if l.Meta.Direction == LINE_SCHEMA_DIRECTION_OUT && k == "dst" {
				continue
			}
			v, ok := m[k]
			if ok {
				if k == "type" && v == "string" {
					continue // 字符串类型,默认不写
				}
				if v == "true" {
					kvArr = append(kvArr, k)
				} else {
					kvArr = append(kvArr, fmt.Sprintf("%s=%s", k, v))
				}
			}
		}
		line := strings.Join(kvArr, ",")
		lineArr = append(lineArr, line)
	}
	out := strings.Join(lineArr, EOF)
	return out
}

type KVpair struct {
	Key   string
	Value string
}

func (kvPair KVpair) Order() int {

	for k, v := range jsonschemalineItemOrder {
		if kvPair.Key == v {
			return k
		}
	}
	return 0
}

type TagLineKVpair []KVpair

func (a TagLineKVpair) Len() int           { return len(a) }
func (a TagLineKVpair) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TagLineKVpair) Less(i, j int) bool { return a[i].Order() < a[j].Order() }

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

	for _, item := range jsonschemaline.Items {
		str := strings.ReplaceAll(item.Fullname, "[]", ".#")
		srcOrDst := fmt.Sprintf("%s.%s", jsonschemaline.Meta.ID, str)
		if item.Src == "" {
			item.Src = srcOrDst
		}
		if item.Dst == "" {
			item.Dst = srcOrDst
		}
	}

	return jsonschemaline, nil
}

// ParseJsonschemalineRaw 解析 jsonschemaline 一行数据
func ParseJsonschemalineRaw(jsonschemalineRaw string) (meta *Meta, item *JsonschemalineItem, err error) {
	jsonschemalineRaw = PretreatJsonschemalineRaw(jsonschemalineRaw)
	kvStrArr := SplitOnUnescapedCommas(jsonschemalineRaw)
	kvMap := make(map[string]string)
	enumList := make([]string, 0)
	constList := make([]string, 0)
	tagLineKVPair := make(TagLineKVpair, 0)
	for _, kvStr := range kvStrArr {
		kvPair := strings.SplitN(kvStr, "=", 2)
		if len(kvPair) == 2 {
			k, v := strings.TrimSpace(kvPair[0]), strings.TrimSpace(kvPair[1])
			switch k {
			case "enum":
				enumList = append(enumList, v)
			case "const":
				constList = append(constList, v)
			default:
				kvMap[k] = v
			}
			tagLineKVPair = append(tagLineKVPair, KVpair{Key: k, Value: v})
		}
	}
	if len(enumList) > 0 {
		jb, err := json.Marshal(enumList)
		if err != nil {
			return nil, nil, err
		}
		kvMap["enum"] = string(jb)
	}
	if len(constList) > 0 {
		jb, err := json.Marshal(constList)
		if err != nil {
			return nil, nil, err
		}
		kvMap["const"] = string(jb)
	}
	jb, err := json.Marshal(kvMap)
	if err != nil {
		return nil, nil, err
	}
	if IsMetaLine(tagLineKVPair) {
		meta = new(Meta)
		err = json.Unmarshal(jb, meta)
		if err != nil {
			return nil, nil, err
		}
		if meta.Version == "" || meta.ID == "" {
			err := errors.Errorf("meta line required version、id ,got:%s", jsonschemalineRaw)
			return nil, nil, err
		}
		switch meta.Direction {
		case LINE_SCHEMA_DIRECTION_IN, LINE_SCHEMA_DIRECTION_OUT:

		default:
			err := errors.Errorf("meta direction must one of  [%s,%s] ,got:%s", LINE_SCHEMA_DIRECTION_IN, LINE_SCHEMA_DIRECTION_OUT, jsonschemalineRaw)
			return nil, nil, err

		}
		return meta, nil, nil
	}
	item = new(JsonschemalineItem)
	err = json.Unmarshal(jb, item)
	if err != nil {
		return nil, nil, err
	}
	if item.Src == "" && item.Dst == "" {
		err := errors.Errorf("at least one of dst/src required ,got :%s", jsonschemalineRaw)
		return nil, nil, err
	}

	item.TagLineKVpair = tagLineKVPair
	return nil, item, nil
}

func (l *Jsonschemaline) JsonSchema() (jsonschemaByte []byte, err error) {
	jsonSchema := new(Schema)
	jsonSchema.Raw2Schema(*l)
	jsonschemaByte, err = jsonSchema.MarshalJSON()
	if err != nil {
		return
	}
	return jsonschemaByte, nil
}
func (l *Jsonschemaline) Jsonschemaline2json() (jsonStr string, err error) {
	jsonStr = ""
	for _, item := range l.Items {
		key := strings.ReplaceAll(item.Fullname, "[]", ".0")
		var value interface{}
		if item.Example != "" {
			value = item.Example
		} else if item.Default != "" {
			value = item.Default
		} else {
			switch item.Type {
			case "int", "integer":
				value = 0
			case "number":
				value = "0"
			case "string":
				value = ""
			}
		}

		jsonStr, err = sjson.Set(jsonStr, key, value)
		if err != nil {
			return "", err
		}
	}
	return jsonStr, nil
}

//PretreatJsonschemalineRaw 处理enum []格式
func PretreatJsonschemalineRaw(tag string) (formatTag string) {
	preg := "enum=\\[(.*)\\]"
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
	formatTag = strings.Trim(formatTag, ",") // 删除前后分号
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

func JsonSchema2LineSchema(jsonschemaStr string) (lineschemaStr string, err error) {
	var jschema Schema
	err = jschema.UnmarshalJSON([]byte(jsonschemaStr))
	if err != nil {
		return "", err
	}
	return jschema.Lineschema()
}

//Json2gsonPatch
func Json2lineSchema(jsonStr string) (out *Jsonschemaline, err error) {
	out = &Jsonschemaline{
		Meta: &Meta{
			Version:   "http://json-schema.org/draft-07/schema#",
			ID:        "example",
			Direction: LINE_SCHEMA_DIRECTION_IN,
		},
		Items: make([]*JsonschemalineItem, 0),
	}
	var input interface{}
	err = json.Unmarshal([]byte(jsonStr), &input)
	if err != nil {
		return nil, err
	}
	rv := reflect.Indirect(reflect.ValueOf(input))
	out.Items = parseOneJsonKey2Line(rv, "")
	return out, nil
}

func parseOneJsonKey2Line(rv reflect.Value, fullname string) (items []*JsonschemalineItem) {
	items = make([]*JsonschemalineItem, 0)
	if rv.IsZero() {
		return items
	}
	rv = reflect.Indirect(rv)
	kind := rv.Kind()
	switch kind {
	case reflect.Int, reflect.Float64, reflect.Int64:
		item := &JsonschemalineItem{
			Type:     "string",
			Format:   "number",
			Fullname: fullname,
			Example:  rv.String(),
		}
		items = append(items, item)
	case reflect.String:
		item := &JsonschemalineItem{
			Type:     "string",
			Fullname: fullname,
			Example:  rv.String(),
		}
		items = append(items, item)
	case reflect.Array, reflect.Slice:
		l := rv.Len()
		for i := 0; i < l; i++ {
			v := rv.Index(i)
			subFullname := fmt.Sprintf("%s[]", fullname)
			subItems := parseOneJsonKey2Line(v, subFullname)
			items = append(items, subItems...)
		}
	case reflect.Map:
		iter := rv.MapRange()
		for iter.Next() {
			k := iter.Key().String()
			subFullname := k
			if fullname != "" {
				subFullname = fmt.Sprintf("%s.%s", fullname, k)
			}
			subItems := parseOneJsonKey2Line(iter.Value(), subFullname)
			items = append(items, subItems...)
		}
	case reflect.Interface, reflect.Ptr:
		rv = rv.Elem()
		return parseOneJsonKey2Line(rv, fullname)
	}
	return items
}
