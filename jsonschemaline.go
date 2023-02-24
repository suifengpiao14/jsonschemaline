package jsonschemaline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/kvstruct"
	"github.com/tidwall/gjson"
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
	ContentEncoding  string       `json:"contentEncoding,omitempty"`   // section 8.3
	ContentMediaType string       `json:"contentMediaType,omitempty"`  // section 8.4
	Title            string       `json:"title,omitempty"`             // section 9.1
	Description      string       `json:"description,omitempty"`       // section 9.1
	Default          string       `json:"default,omitempty"`           // section 9.2
	Deprecated       bool         `json:"deprecated,omitempty,string"` // section 9.3
	ReadOnly         bool         `json:"readOnly,omitempty,string"`   // section 9.4
	WriteOnly        bool         `json:"writeOnly,omitempty,string"`  // section 9.4
	Example          string       `json:"example,omitempty"`           // section 9.5
	Src              string       `json:"src,omitempty"`
	Dst              string       `json:"dst,omitempty"`
	Fullname         string       `json:"fullname,omitempty"`
	TagLineKVpair    kvstruct.KVS `json:"-"`
}

func (jItem JsonschemalineItem) Json() (jsonStr string) {
	b, _ := json.Marshal(jItem)
	jsonStr = string(b)
	return jsonStr
}

func (jItem JsonschemalineItem) jsonSchemaItem() (jsonStr string) {
	copy := jItem
	copy.Required = false // 转换成json schema时 required 单独处理
	// 这部分字段隐藏
	copy.Fullname = ""
	copy.Dst = ""
	copy.Src = ""
	b, _ := json.Marshal(copy)
	jsonStr = string(b)
	return jsonStr
}

func (jItem JsonschemalineItem) ToJsonSchemaKVS() (kvs kvstruct.KVS, err error) {
	kvs = make(kvstruct.KVS, 0)
	arrSuffix := "[]"
	fullname := strings.Trim(jItem.Fullname, ".")
	if !strings.HasPrefix(fullname, arrSuffix) {
		fullname = fmt.Sprintf(".%s", fullname) //增加顶级对象
	}
	arr := strings.Split(fullname, ".")
	kv := kvstruct.KV{
		Key:   `$schema`,
		Value: `http://json-schema.org/draft-07/schema#`,
	}
	kvs = append(kvs, kv)
	prefix := ""
	l := len(arr)
	for i := 0; i < l; i++ {
		key := arr[i]
		//处理数组
		if strings.HasSuffix(key, arrSuffix) {
			key = strings.TrimSuffix(key, arrSuffix)
			prefix = strings.Trim(fmt.Sprintf("%s.%s", prefix, key), ".")
			kv := kvstruct.KV{
				Key:   strings.Trim(fmt.Sprintf("%s.type", prefix), "."),
				Value: "array",
			}
			kvs = append(kvs, kv)
			if i == l-1 {
				kv = kvstruct.KV{
					Key:   strings.Trim(fmt.Sprintf("%s.items", prefix), "."),
					Value: jItem.jsonSchemaItem(),
				}
				continue
			}
			prefix = fmt.Sprintf("%s.items", prefix)
			kv = kvstruct.KV{
				Key:   strings.Trim(fmt.Sprintf("%s.type", prefix), "."),
				Value: "object",
			}
			kvs = append(kvs, kv)
			prefix = fmt.Sprintf("%s.properties", prefix)
			continue
		}

		//处理对象
		if i == l-1 {
			if jItem.Required {
				parentKey := strings.TrimSuffix(prefix, ".properties")
				kv := kvstruct.KV{
					Key:   strings.Trim(fmt.Sprintf("%s.required.-1", parentKey), "."),
					Value: key,
				}
				kvs = append(kvs, kv)
			}
			kv := kvstruct.KV{
				Key:   strings.Trim(fmt.Sprintf("%s.%s", prefix, key), "."),
				Value: jItem.jsonSchemaItem(),
			}
			kvs = append(kvs, kv)
			continue
		}

		prefix = strings.Trim(fmt.Sprintf("%s.%s", prefix, key), ".")
		kv := kvstruct.KV{
			Key:   strings.Trim(fmt.Sprintf("%s.type", prefix), "."),
			Value: "object",
		}
		kvs = append(kvs, kv)
		prefix = fmt.Sprintf("%s.properties", prefix)
	}
	return kvs, nil
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

func IsMetaLine(lineTags kvstruct.KVS) bool {
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
	lineArr = append(lineArr, fmt.Sprintf("version=%s,direction=%s,id=%s", l.Meta.Version, l.Meta.Direction, l.Meta.ID))
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

// ParseMultiJsonSchemaline 解析多个 jsonschemaline
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

// ParseJsonschemaline 解析单个 jsonschemaline
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
	tagLineKVPair := make(kvstruct.KVS, 0)
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
			tagLineKVPair = append(tagLineKVPair, kvstruct.KV{Key: k, Value: v})
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
	kvs := kvstruct.KVS{
		{Key: "$schema", Value: "http://json-schema.org/draft-07/schema#"},
	}
	for _, item := range l.Items {
		subKvs, err := item.ToJsonSchemaKVS()
		if err != nil {
			return nil, err
		}
		kvs.Add(subKvs...)
	}

	jsonschemaByte = []byte("")
	for _, kv := range kvs {
		if kvstruct.IsJsonStr(kv.Value) {
			jsonschemaByte, err = sjson.SetRawBytes(jsonschemaByte, kv.Key, []byte(kv.Value))
			if err != nil {
				return nil, err
			}
			continue
		}
		jsonschemaByte, err = sjson.SetBytes(jsonschemaByte, kv.Key, kv.Value)
		if err != nil {
			return nil, err
		}
	}
	return jsonschemaByte, nil
}
func (l *Jsonschemaline) JsonExample() (jsonExample string, err error) {
	jsonExample = ""
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
		existsResult := gjson.Get(jsonExample, key)
		if existsResult.IsArray() || existsResult.IsObject() {
			continue
		}
		jsonExample, err = sjson.Set(jsonExample, key, value)
		if err != nil {
			return "", err
		}
	}
	return jsonExample, nil
}

// Deprecated see JsonExample
func (l *Jsonschemaline) Jsonschemaline2json() (jsonStr string, err error) {
	return l.JsonExample()
}

type Struct struct {
	IsRoot     bool
	Name       string
	Lineschema string
	Attrs      []StructAttr
}

// AddAttrIgnore 已经存在则跳过
func (s *Struct) AddAttrIgnore(attrs ...StructAttr) {
	if len(s.Attrs) == 0 {
		s.Attrs = make([]StructAttr, 0)
	}
	for _, attr := range attrs {
		if _, exists := s.GetAttr(attr.Name); exists {
			continue
		}
		s.Attrs = append(s.Attrs, attr)
	}

}
func (s *Struct) GetAttr(attrName string) (structAttr *StructAttr, exists bool) {
	for _, attr := range s.Attrs {
		if attr.Name == attrName {
			return &attr, true
		}
	}
	return nil, false
}

type StructAttr struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

type Structs []*Struct

func (s *Structs) Json() (str string) {
	b, _ := json.Marshal(s)
	str = string(b)
	return str
}

func (s *Structs) GetRoot() (struc *Struct, exists bool) {
	for _, stru := range *s {
		if stru.IsRoot {
			return stru, true
		}
	}
	return struc, false
}

func (s *Structs) Get(name string) (struc *Struct, exists bool) {
	for _, stru := range *s {
		if stru.Name == name {
			return stru, true
		}
	}
	return struc, false
}

func (s *Structs) AddIngore(structs ...*Struct) {
	if len(*s) == 0 {
		*s = make(Structs, 0)
	}
	for _, structRef := range structs {
		if _, exists := s.Get(structRef.Name); exists {
			continue
		}
		*s = append(*s, structRef)

	}
}

func (l *Jsonschemaline) ToSturct() (structs Structs) {
	arraySuffix := "[]"
	structs = make(Structs, 0)
	id := string(l.Meta.ID)
	rootStructName := ToCamel(id)
	rootStruct := &Struct{
		IsRoot:     true,
		Name:       rootStructName,
		Attrs:      make([]StructAttr, 0),
		Lineschema: l.String(),
	}
	structs.AddIngore(rootStruct)
	for _, item := range l.Items {
		if item.Fullname == "" {
			continue
		}
		withRootFullname := strings.Trim(fmt.Sprintf("%s.%s", id, item.Fullname), ".")
		nameArr := strings.Split(withRootFullname, ".")
		nameCount := len(nameArr)
		for i := 1; i < nameCount; i++ { //i从1开始,0 为root,已处理
			baseName := nameArr[i]
			if i < nameCount-1 { // 非最后一个,即为上级的attr,又为下级的struct
				realBaseName := strings.TrimSuffix(baseName, arraySuffix)
				isArray := baseName != realBaseName
				attrName := ToCamel(realBaseName)
				subStructName := ToCamel(strings.Join(nameArr[:i+1], "_"))
				if _, ok := structs.Get(subStructName); ok { // 存在跳过
					continue
				}
				//不存在,新增父类属性,增加子类
				typ := subStructName
				if isArray {
					typ = fmt.Sprintf("[]%s", typ)
				}
				attr := StructAttr{
					Name: attrName,
					Type: typ,
					Tag:  fmt.Sprintf(`json:"%s"`, ToLowerCamel(attrName)),
				}
				parentStructName := ToCamel(strings.Join(nameArr[:i], "_"))
				parentStruct, _ := structs.Get(parentStructName) // 一定存在
				parentStruct.AddAttrIgnore(attr)
				subStruct := &Struct{
					IsRoot: false,
					Name:   subStructName,
				}
				structs.AddIngore(subStruct)
				continue
			}

			// 最后一个
			realBaseName := strings.TrimSuffix(baseName, arraySuffix)
			isArray := baseName != realBaseName
			attrName := ToCamel(realBaseName)
			parentStructName := ToCamel(strings.Join(nameArr[:i], "_"))
			parentStruct, _ := structs.Get(parentStructName)

			typ := item.Type
			tag := fmt.Sprintf(`json:"%s"`, ToLowerCamel(attrName))
			if typ == "string" {
				switch item.Format {
				case "int", "number":
					typ = "int"
					tag = fmt.Sprintf(`json:"%s,string"`, ToLowerCamel(attrName))
				}
			}
			if l.Meta.Direction == LINE_SCHEMA_DIRECTION_IN && !item.Required { //当作入参时,非必填字断,使用引用
				typ = fmt.Sprintf("*%s", typ)
			}
			if isArray {
				typ = fmt.Sprintf("[]%s", typ)
			}

			newAttr := &StructAttr{
				Name:    ToCamel(attrName),
				Type:    typ,
				Tag:     tag,
				Comment: item.Comments,
			}
			attr, ok := parentStruct.GetAttr(attrName)
			if ok { //已经存在,修正类型和备注
				typ := newAttr.Type
				if strings.HasPrefix(attr.Type, "[]") && !strings.HasPrefix(typ, "[]") {
					typ = fmt.Sprintf("[]%s", typ)
				}
				attr.Type = typ
				if newAttr.Comment != "" {
					attr.Comment = newAttr.Comment
				}
				continue
			}
			// 不存在,新增
			parentStruct.AddAttrIgnore(*newAttr)
		}
	}

	return structs
}

func (l *Jsonschemaline) GjsonPath(formatPath func(format string, src string, item *JsonschemalineItem) (path string)) (gjsonPath string) {
	m := &map[string]interface{}{}
	for _, item := range l.Items {
		dst, src, format := item.Dst, item.Src, item.Format
		if formatPath != nil {
			src = formatPath(format, src, item)
		}
		dst = strings.ReplaceAll(dst, ".#", "[]") //替换成[],方便后续遍历
		arr := strings.Split(dst, ".")
		l := len(arr)
		var ref = new(map[string]interface{})
		*ref = *m
		for i, key := range arr {
			if l == i+1 {
				(*ref)[key] = src
			}
			if _, ok := (*ref)[key]; !ok {
				temp := &map[string]interface{}{}
				(*ref)[key] = temp
			}
			tmp, ok := (*ref)[key].(*map[string]interface{}) //递进
			if ok {
				*ref = *tmp
			}
		}

	}
	w := recursionWrite(m)
	gjsonPath = fmt.Sprintf("{%s}", w.String())
	return gjsonPath
}

// 生成路径
func recursionWrite(m *map[string]interface{}) (w bytes.Buffer) {
	writeComma := false
	for k, v := range *m {
		if writeComma {
			w.WriteString(",")
		}
		writeComma = true
		ref, ok := v.(*map[string]interface{})
		if !ok {
			w.WriteString(fmt.Sprintf("%s:%s", k, v))
			continue
		}
		subw := recursionWrite(ref)
		if strings.HasSuffix(k, "[]") {
			k = strings.TrimRight(k, "[]")
			subStr := fmt.Sprintf("%s:{%s}|@group", k, subw.String())
			w.WriteString(subStr)
		} else {
			subStr := fmt.Sprintf("%s:{%s}", k, subw.String())
			w.WriteString(subStr)
		}
	}
	return w
}

// PretreatJsonschemalineRaw 处理enum []格式
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

// Json2gsonPatch
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
