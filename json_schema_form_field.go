package jsonschemaline

type JsonSchemaFormFiledOut struct {
	Fullname    string                     `json:"fullname"`    // 名称/全称
	Type        string                     `json:"type"`        // 类型 枚举值
	Title       string                     `json:"title"`       // 标题
	Description string                     `json:"description"` // 描述
	Format      string                     `json:"format"`      // 类型格式
	More        JsonSchemaFormFiledOutMore `json:"more"`
}

type JsonSchemaFormFiledOutMore struct {
	Enum             []string `json:"enum"`             // 枚举值
	EnumNames        []string `json:"enumNames"`        // 枚举值标题
	Comment          string   `json:"comment"`          // 备注
	Const            string   `json:"const"`            // 常量
	MultipleOf       int      `json:"multipleOf"`       // 多值
	Maximum          int      `json:"maximum"`          // 最大值
	ExclusiveMaximum bool     `json:"exclusiveMaximum"` // 是否包含最大值
	Minimum          int      `json:"minimum"`          // 最小值
	ExclusiveMinimum bool     `json:"exclusiveMinimum"` // 是否包含最小值
	MaxLength        int      `json:"maxLength"`        // 最大长度
	MinLength        int      `json:"minLength"`        // 最小长度
	Pattern          string   `json:"pattern"`          // 匹配格式
	MaxItems         int      `json:"maxItems"`         // 最大项数
	MinItems         int      `json:"minItems"`         // 最小项数
	UniqueItems      bool     `json:"uniqueItems"`      // 数组唯一
	MaxContains      uint     `json:"maxContains"`      // 符合Contains规则最大数量
	MinContains      uint     `json:"minContains"`      // 符合Contains规则最小数量
	MaxProperties    int      `json:"maxProperties"`    // 对象最多属性个数
	MinProperties    int      `json:"minProperties"`    // 对象最少属性个数
	Required         bool     `json:"required"`         // 是否必须
	ContentEncoding  string   `json:"contentEncoding"`  // 内容编码
	ContentMediaType string   `json:"contentMediaType"` // 内容格式
	Default          string   `json:"default"`          // 默认值
	Deprecated       bool     `json:"deprecated"`       // 是否弃用
	ReadOnly         bool     `json:"readOnly"`         // 只读
	WriteOnly        bool     `json:"writeOnly"`        // 只写
	Example          string   `json:"example"`          // 案例
	Examples         string   `json:"examples"`         // 案例集合
	AllowEmptyValue  bool     `json:"allowEmptyValue"`  // 是否可以为空
}

//GetJsonSchemaFormFileds jsonschema 表单字段
func GetJsonSchemaFormFileds() (formFields []JsonSchemaFormFiledOut) {
	formFields = []JsonSchemaFormFiledOut{
		{Fullname: "fullname", Type: "string", Title: "名称/全称", Description: "名称/全称"},
		{Fullname: "type", Type: "string", Title: "类型", Description: "类型", More: JsonSchemaFormFiledOutMore{Enum: []string{"int", "float", "string"}, EnumNames: []string{"整型", "浮点型", "字符串"}}},
		{Fullname: "enum", Type: "string", Title: "枚举值", Description: "枚举值"},
		{Fullname: "enumNames", Type: "string", Title: "枚举值标题", Description: "枚举值标题"},
		{Fullname: "comment", Type: "string", Title: "备注", Description: "备注"},
		{Fullname: "const", Type: "string", Title: "常量", Description: "常量"},
		{Fullname: "multipleOf", Type: "string", Format: "int", Title: "多值", Description: "多值"},
		{Fullname: "maximum", Type: "string", Format: "int", Title: "最大值", Description: "最大值"},
		{Fullname: "exclusiveMaximum", Type: "string", Format: "bool", Title: "是否包含最大值", Description: "是否包含最大值"},
		{Fullname: "minimum", Type: "string", Format: "int", Title: "最小值", Description: "最小值"},
		{Fullname: "exclusiveMinimum", Type: "string", Format: "bool", Title: "是否包含最小值", Description: "是否包含最小值"},
		{Fullname: "maxLength", Type: "string", Format: "int", Title: "最大长度", Description: "最大长度"},
		{Fullname: "minLength", Type: "string", Format: "int", Title: "最小长度", Description: "最小长度"},
		{Fullname: "pattern", Type: "string", Title: "匹配格式", Description: "匹配格式"},
		{Fullname: "maxItems", Type: "string", Format: "int", Title: "最大项数", Description: "最大项数"},
		{Fullname: "minItems", Type: "string", Format: "int", Title: "最小项数", Description: "最小项数"},
		{Fullname: "uniqueItems", Type: "string", Format: "bool", Title: "数组唯一", Description: "数组唯一"},
		{Fullname: "maxContains", Type: "string", Format: "uint", Title: "符合Contains规则最大数量", Description: "符合Contains规则最大数量"},
		{Fullname: "minContains", Type: "string", Format: "uint", Title: "符合Contains规则最小数量", Description: "符合Contains规则最小数量"},
		{Fullname: "maxProperties", Type: "string", Format: "int", Title: "对象最多属性个数", Description: "对象最多属性个数"},
		{Fullname: "minProperties", Type: "string", Format: "int", Title: "对象最少属性个数", Description: "对象最少属性个数"},
		{Fullname: "required", Type: "string", Format: "bool", Title: "是否必须", Description: "是否必须"},
		{Fullname: "format", Type: "string", Title: "类型格式", Description: "类型格式"},
		{Fullname: "contentEncoding", Type: "string", Title: "内容编码", Description: "内容编码"},
		{Fullname: "contentMediaType", Type: "string", Title: "内容格式", Description: "内容格式"},
		{Fullname: "title", Type: "string", Title: "标题", Description: "标题"},
		{Fullname: "description", Type: "string", Title: "描述", Description: "描述"},
		{Fullname: "default", Type: "string", Title: "默认值", Description: "默认值"},
		{Fullname: "deprecated", Type: "string", Format: "bool", Title: "是否弃用", Description: "是否弃用"},
		{Fullname: "readOnly", Type: "string", Format: "bool", Title: "只读", Description: "只读"},
		{Fullname: "writeOnly", Type: "string", Format: "bool", Title: "只写", Description: "只写"},
		{Fullname: "example", Type: "string", Title: "案例", Description: "案例"},
		{Fullname: "examples", Type: "string", Title: "案例集合", Description: "案例集合"},
		{Fullname: "allowEmptyValue", Type: "string", Format: "bool", Title: "是否可以为空", Description: "是否可以为空"},
	}
	return formFields
}
