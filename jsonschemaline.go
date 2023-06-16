package jsonschemaline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/suifengpiao14/funcs"
	_ "github.com/suifengpiao14/gjsonmodifier"
	"github.com/suifengpiao14/kvstruct"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type JsonschemalineItem struct {
	Comments string `json:"comment,omitempty"` // section 8.3

	Type             string `json:"type,omitempty"`                    // section 6.1.1
	Enum             string `json:"enum,omitempty"`                    // section 6.1.2
	EnumNames        string `json:"enumNames,omitempty"`               // section 6.1.2
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
	Examples         string       `json:"examples,omitempty"`          // section 9.5
	Src              string       `json:"src,omitempty"`
	Dst              string       `json:"dst,omitempty"`
	Fullname         string       `json:"fullname,omitempty"`
	AllowEmptyValue  bool         `json:"allowEmptyValue,omitempty,string"`
	TagLineKVpair    kvstruct.KVS `json:"-"`
}

func (jItem JsonschemalineItem) String() (jsonStr string) {
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

func (jItem JsonschemalineItem) ToKVS(namespance string) (kvs kvstruct.KVS) {
	jsonStr := jItem.String()
	kvs = kvstruct.JsonToKVS(jsonStr, namespance)
	return kvs
}
func (jItem JsonschemalineItem) enum2Array() (enum []string, enumNames []string, err error) {
	if jItem.Enum != "" {
		var enumI []interface{}
		err = json.Unmarshal([]byte(jItem.Enum), &enumI)
		if err != nil {
			return nil, nil, err
		}
		enum = make([]string, 0)
		for _, e := range enumI {
			enum = append(enum, cast.ToString(e))
		}
	}
	if jItem.EnumNames != "" {
		var enumNamesI []interface{}
		err = json.Unmarshal([]byte(jItem.EnumNames), &enumNamesI)
		if err != nil {
			return nil, nil, err
		}
		enumNames = make([]string, 0)
		for _, e := range enumNamesI {
			enumNames = append(enumNames, cast.ToString(e))
		}
	}
	return enum, enumNames, nil
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
				fullKey := strings.Trim(fmt.Sprintf("%s.items", prefix), ".")
				attrKvs := jItem.ToKVS(fullKey)
				kvs.AddReplace(attrKvs...)
				enum, enumNames, err := jItem.enum2Array()
				if err != nil {
					return nil, err
				}
				subKvs := enumNames2KVS(enum, enumNames, fullKey)
				kvs.AddReplace(subKvs...)
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
				kvs.AddReplace(kv)
			}
			fullKey := strings.Trim(fmt.Sprintf("%s.%s", prefix, key), ".")
			attrKvs := jItem.ToKVS(fullKey)
			kvs.AddReplace(attrKvs...)
			enum, enumNames, err := jItem.enum2Array()
			if err != nil {
				return nil, err
			}
			subKvs := enumNames2KVS(enum, enumNames, fullKey)
			kvs.AddReplace(subKvs...)
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

func enumNames2KVS(enum []string, enumNames []string, prefix string) (kvs kvstruct.KVS) {
	kvs = make(kvstruct.KVS, 0)
	if len(enumNames) < 1 {
		return kvs
	}
	enumLen := len(enum)
	for i, enumName := range enumNames {
		if i >= enumLen {
			continue
		}
		enum := enum[i]
		kv := kvstruct.KV{
			Key:   strings.Trim(fmt.Sprintf("%s.oneOf.%d.const", prefix, i), "."),
			Value: enum,
		}
		kvs.Add(kv)
		kv = kvstruct.KV{
			Key:   strings.Trim(fmt.Sprintf("%s.oneOf.%d.title", prefix, i), "."),
			Value: enumName,
		}
		kvs.Add(kv)
	}
	return kvs
}

var jsonschemalineItemOrder = []string{
	"fullname", "src", "dst", "type", "format", "pattern", "enum", "required", "allowEmptyValue", "title", "description", "default", "comment", "example", "deprecated", "const",
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
		if gjson.GetBytes(jsonschemaByte, kv.Key).Exists() { // 已经存在的，不覆盖（防止 array、object 在其子属性说明后，导致覆盖）
			continue
		}
		if kvstruct.IsJsonStr(kv.Value) {
			jsonschemaByte, err = sjson.SetRawBytes(jsonschemaByte, kv.Key, []byte(kv.Value))
			if err != nil {
				return nil, err
			}
			continue
		}
		var value interface{}
		value = kv.Value
		baseKey := BaseName(kv.Key)
		switch baseKey {
		case "exclusiveMaximum", "exclusiveMinimum", "deprecated", "readOnly", "writeOnly", "uniqueItems":
			value = kv.Value == "true"
		case "multipleOf", "maximum", "minimum", "maxLength", "minLength", "maxItems", "minItems", "maxContains", "minContains", "maxProperties", "minProperties":
			value, _ = strconv.Atoi(kv.Value)
		}
		jsonschemaByte, err = sjson.SetBytes(jsonschemaByte, kv.Key, value)
		if err != nil {
			return nil, err
		}
	}
	return jsonschemaByte, nil
}
func ReplacePathSpecalChar(path string) (newPath string) {
	replacer := strings.NewReplacer("|", "\\|", "#", "\\#", "@", "\\@", "*", "\\*", "?", "\\?")
	return replacer.Replace(path)
}
func (l *Jsonschemaline) JsonExample() (jsonExample string, err error) {
	jsonExample = ""
	for _, item := range l.Items {
		key := strings.ReplaceAll(item.Fullname, "[]", ".0")
		var value interface{}
		if item.Examples != "" {
			value = item.Examples
		} else if item.Example != "" {
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
		key = ReplacePathSpecalChar(key)
		existsResult := gjson.Get(jsonExample, key)
		if existsResult.IsArray() || existsResult.IsObject() { //支持array、object 整体设置example
			if str, ok := value.(string); ok {
				jsonExample, err = sjson.SetRaw(jsonExample, key, str)
				if err != nil {
					return "", err
				}
			}
			continue
		}
		jsonExample, err = sjson.Set(jsonExample, key, value)
		if err != nil {
			return "", err
		}
	}
	return jsonExample, nil
}

type DefaultJson struct {
	ID      string
	Version string
	Json    string
}

func (l *Jsonschemaline) DefaultJson() (defaultJson *DefaultJson, err error) {
	defaultJson = new(DefaultJson)
	id := l.Meta.ID.String()
	defaultJson.ID = id
	defaultJson.Version = l.Meta.Version
	kvmap := make(map[string]string)
	for _, item := range l.Items {
		if item.Default != "" || item.AllowEmptyValue {
			path := strings.ReplaceAll(item.Fullname, "[]", ".#")
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

func (l *Jsonschemaline) ToSturct() (structs Structs) {
	arraySuffix := "[]"
	structs = make(Structs, 0)
	id := string(l.Meta.ID)
	rootStructName := funcs.ToCamel(id)
	rootStruct := &Struct{
		IsRoot:     true,
		Name:       rootStructName,
		Attrs:      make([]*StructAttr, 0),
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
			parentStructName := funcs.ToCamel(strings.Join(nameArr[:i], "_"))
			parentStruct, _ := structs.Get(parentStructName) // 一定存在
			baseName := nameArr[i]
			realBaseName := strings.TrimSuffix(baseName, arraySuffix)
			isArray := baseName != realBaseName || item.Type == "array"
			attrName := funcs.ToCamel(realBaseName)
			comment := item.Comments
			if comment == "" {
				comment = item.Description
			}
			if i < nameCount-1 { // 非最后一个,即为上级的attr,又为下级的struct
				subStructName := funcs.ToCamel(strings.Join(nameArr[:i+1], "_"))
				attrType := subStructName
				if isArray {
					attrType = fmt.Sprintf("[]%s", attrType)
				}
				attr := StructAttr{
					Name: attrName,
					Type: attrType,
					Tag:  fmt.Sprintf(`json:"%s"`, funcs.ToLowerCamel(attrName)),
					//Comment: comment,// 符合类型comment 无意义，不增加
				}
				parentStruct.AddAttrReplace(attr)
				subStruct := &Struct{
					IsRoot: false,
					Name:   subStructName,
				}
				structs.AddIngore(subStruct)
				continue
			}
			format := item.Format

			switch format { // 格式化format
			case "number":
				format = "int"
			case "float":
				format = "float64"
			}

			// 最后一个
			typ := item.Type
			if typ == "string" && format != "" {
				typ = format
			}
			tag := fmt.Sprintf(`json:"%s"`, funcs.ToLowerCamel(attrName))
			if l.Meta.Direction == LINE_SCHEMA_DIRECTION_IN && !item.Required { //当作入参时,非必填字断,使用引用
				typ = fmt.Sprintf("*%s", typ)
			}
			if isArray {
				if typ == "array" {
					typ = "interface{}"
					if format != "" {
						typ = format
					}
				}
				typ = fmt.Sprintf("[]%s", typ)
			}

			newAttr := &StructAttr{
				Name:    funcs.ToCamel(attrName),
				Type:    typ,
				Tag:     tag,
				Comment: comment,
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

// GjsonPathWithDefaultFormat 生成格式化的jsonpath，用来重新格式化数据,比如入参字段类型全为字符串，在format中标记了实际类型，可以通过该方法获取转换数据的gjson path，从入参中提取数据后，对应字段类型就以format为准，此处仅仅提供有创意的案例，更多可以依据该思路扩展
func (l *Jsonschemaline) GjsonPathWithDefaultFormat(ignoreID bool) (gjsonPath string) {
	switch l.Meta.Direction {
	case LINE_SCHEMA_DIRECTION_IN:
		gjsonPath = l.GjsonPath(ignoreID, FormatPathFnByFormatIn)
	case LINE_SCHEMA_DIRECTION_OUT:
		gjsonPath = l.GjsonPath(ignoreID, FormatPathFnByFormatOut)
	}

	return gjsonPath
}

func (l *Jsonschemaline) GjsonPath(ignoreID bool, formatPath func(format string, src string, item *JsonschemalineItem) (path string)) (gjsonPath string) {
	m := &map[string]interface{}{}
	for _, item := range l.Items {
		switch strings.ToLower(item.Type) {
		case "array", "object": // 数组、对象需要遍历内部结构,忽略外部的path
			continue
		}
		dst, src, format := item.Dst, item.Src, item.Format
		dst = strings.ReplaceAll(dst, ".#", "[]") //替换成[],方便后续遍历
		if ignoreID {
			switch l.Meta.Direction {
			case LINE_SCHEMA_DIRECTION_IN:
				src = strings.TrimPrefix(src, fmt.Sprintf("%s.", l.Meta.ID))
			case LINE_SCHEMA_DIRECTION_OUT:
				dst = strings.TrimPrefix(dst, fmt.Sprintf("%s.", l.Meta.ID))
			}

		}
		if formatPath != nil {
			src = formatPath(format, src, item)
		}
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

// 使用format 属性格式化转换后的路径
func FormatPathFnByFormatIn(format string, src string, item *JsonschemalineItem) (path string) {
	path = src
	switch format {
	case "int", "float", "number":
		path = fmt.Sprintf("%s.@tonum", src)
	case "bool":
		path = fmt.Sprintf("%s.@tobool", src)
	}
	return path
}

// 使用format 属性格式化转换后的路径
func FormatPathFnByFormatOut(format string, src string, item *JsonschemalineItem) (path string) {
	path = src
	if item.Type == "string" {
		path = fmt.Sprintf("%s.@tostring", src)
	}
	return path
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
			k = strings.TrimSuffix(k, "[]")
			w.WriteString(fmt.Sprintf("%s:%s", k, v))
			continue
		}
		subw := recursionWrite(ref)
		isArray := false
		subwKey := subw.String()
		var subStr string
		if strings.HasSuffix(k, "[]") {
			isArray = true
			k = strings.TrimRight(k, "[]")
		}

		if strings.Contains(subwKey, ".#.") {
			isArray = true
		}

		if isArray {
			subStr = fmt.Sprintf("%s:{%s}|@group", k, subwKey)
		} else {
			subStr = fmt.Sprintf("%s:{%s}", k, subwKey)
		}
		w.WriteString(subStr)
	}
	return w
}

// Json2lineSchema
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

const (
	//jsonschema 本身支持的格式描述(后台表单经常用到)
	JsonSchemaLineSchema = `
	version=http://json-schema.org/draft-07/schema#,direction=in,id=jsonschema
	fullname=comment,dst=comment,title=备注
	fullname=type,dst=type,title=类型
	fullname=enum,dst=enum,type=array,format=string,title=枚举值
	fullname=enumNames,dst=enumNames,type=array,format=string,title=枚举值标题
	fullname=const,dst=const,title=常量
	fullname=multipleOf,dst=multipleOf,format=int,title=多值
	fullname=maximum,dst=maximum,format=int,title=最大值
	fullname=exclusiveMaximum,dst=exclusiveMaximum,format=bool,title=是否包含最大值
	fullname=minimum,dst=minimum,format=int,title=最小值
	fullname=exclusiveMinimum,dst=exclusiveMinimum,format=bool,title=是否包含最小值
	fullname=maxLength,dst=maxLength,format=int,title=最大长度
	fullname=minLength,dst=minLength,format=int,title=最小长度
	fullname=pattern,dst=pattern,title=匹配格式
	fullname=maxItems,dst=maxItems,format=int,title=最大项数
	fullname=minItems,dst=minItems,format=int,title=最小项数
	fullname=uniqueItems,dst=uniqueItems,format=bool,title=数组唯一
	fullname=maxContains,dst=maxContains,format=uint,title=符合Contains规则最大数量
	fullname=minContains,dst=minContains,format=uint,title=符合Contains规则最小数量
	fullname=maxProperties,dst=maxProperties,format=int,title=对象最多属性个数
	fullname=minProperties,dst=minProperties,format=int,title=对象最少属性个数
	fullname=required,dst=required,format=bool,title=是否必须
	fullname=format,dst=format,title=类型格式
	fullname=contentEncoding,dst=contentEncoding,title=内容编码
	fullname=contentMediaType,dst=contentMediaType,title=内容格式
	fullname=title,dst=title,title=标题
	fullname=description,dst=description,title=描述
	fullname=default,dst=default,title=默认值
	fullname=deprecated,dst=deprecated,format=bool,title=是否弃用
	fullname=readOnly,dst=readOnly,format=bool,title=只读
	fullname=writeOnly,dst=writeOnly,format=bool,title=只写
	fullname=example,dst=example,title=案例
	fullname=examples,dst=examples,title=案例集合
	fullname=src,dst=src,title=源数据
	fullname=dst,dst=dst,title=目标数据
	fullname=fullname,dst=fullname,title=名称/全称
	fullname=allowEmptyValue,dst=allowEmptyValue,format=bool,title=是否可以为空
	`
)

// GetJsonSchemaScheme 返回jsonschema本身的schema
func GetJsonSchemaSchema() (schema string) {
	lineschema, err := ParseJsonschemaline(JsonSchemaLineSchema)
	if err != nil {
		panic(err)
	}
	b, err := lineschema.JsonSchema()
	if err != nil {
		panic(err)
	}
	schema = string(b)
	return schema
}
