package jsonschemaline_test

import (
	"encoding/json"
	"fmt"
	"testing"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/jsonschemaline"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var jsonStr = `{"$schema":"http://json-schema.org/draft-07/schema#","$id":"execAPIInquiryScreenIdentifyUpdate","type":"object","required":["config"],"properties":{"config":{"properties":{"id":{"type":"string","format":"number"},"status":{"type":"string","enum":["1","2"]},"identify":{"type":"string","minLength":1},"merchantId":{"type":"string","format":"number"},"merchantName":{"type":"string","minLength":1},"operateName":{"type":"string","minLength":1},"storeId":{"type":"string","format":"number"},"storeName":{"type":"string","minLength":1}},"type":"object","required":["id","status","identify","merchantId","merchantName","operateName","storeId","storeName"]}}}`

func TestJsonSchemalineString(t *testing.T) {
	lineschema, err := jsonschemaline.ParseJsonschemaline(schemalineOut)
	if err != nil {
		panic(err)
	}

	fmt.Println(lineschema.String())

}

func TestJsonSchema2LineSchema(t *testing.T) {
	lineschemaStr, err := jsonschemaline.JsonSchema2LineSchema(jsonStr)
	if err != nil {
		panic(err)
	}

	fmt.Println(lineschemaStr)

}

func TestJson2lineSchema(t *testing.T) {
	var jsonStr = `
		{
			"config":{
				"id":"1",
				"status":"2",
				"identify":"abcde",
				"merchantId":"123",
				"merchantName":"测试商户",
				"operateName":"彭政",
				"storeId":"1",
				"storeName":"门店名称"
			}
		}
	`
	lineschema, err := jsonschemaline.Json2lineSchema(jsonStr)
	if err != nil {
		panic(err)
	}
	fmt.Println(lineschema.String())
}

func TestLine2tpl(t *testing.T) {
	line := `
	version=http://json-schema.org/draft-07/schema,id=output,direction=out
fullname=items[].class,src=GetByKeyOut.0.Fvalue.#.class,required
fullname=items[].priceType,src=GetByKeyOut.0.Fvalue.#.priceType,required
fullname=items[].minPrice,src=GetByKeyOut.0.Fvalue.#.minPrice,required
fullname=items[].maxPrice,src=GetByKeyOut.0.Fvalue.#.maxPrice,required
fullname=items[].maxRate,src=GetByKeyOut.0.Fvalue.#.maxRate,required
fullname=items[].minRate,src=GetByKeyOut.0.Fvalue.#.minRate,required
fullname=items[].weight,src=GetByKeyOut.0.Fvalue.#.weight,required
`
	lineschema, err := jsonschemaline.ParseJsonschemaline(line)
	if err != nil {
		panic(err)
	}
	instructTpl := jsonschemaline.ParseInstructTp(*lineschema)
	inputTpl := instructTpl.String()
	fmt.Println(inputTpl)
}

func TestSjson(t *testing.T) {
	val := []string{"1", "2", "3"}
	out, err := sjson.Set("", "output.items.index.-1", val)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
func TestJsonExample(t *testing.T) {

	l := NewLineSchema()
	jsonStr, err := l.JsonExample()
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonStr)
}
func TestJsonExample2(t *testing.T) {
	lineschema := `
version=http://json-schema.org/draft-07/schema#,direction=out,id=example
fullname=index,type=string,dst=index
fullname=position,type=string,dst=position
fullname=httpStatus,type=string,dst=httpStatus
fullname=total,type=string,example=60,dst=total
fullname=type,type=string,dst=type
fullname=description,type=string,dst=description
fullname=name,type=string,dst=name
fullname=size,type=string,example=10,dst=size
`
	l, err := jsonschemaline.ParseJsonschemaline(lineschema)
	if err != nil {
		panic(err)
	}
	jsonStr, err := l.JsonExample()
	if err != nil {
		panic(err)
	}
	expected := `{"index":"","position":"","httpStatus":"","total":"60","type":"","description":"","name":"","size":"10"}`
	ok := jsonpatch.Equal([]byte(expected), []byte(jsonStr))
	assert.Equal(t, true, ok)
}
func TestJsonschemaline2json3(t *testing.T) {
	lineschema := `
	version=http://json-schema.org/draft-07/schema#,direction=out,id=example
	fullname=pagination.index,type=string,example=0,dst=pagination.index
	fullname=pagination.total,type=string,example=60,dst=pagination.total
	fullname=items[].link,type=string,dst=items[].link
	fullname=items[].id,type=string,example=0,dst=items[].id
	fullname=message,type=string,example=-,dst=message
	fullname=pagination.size,type=string,example=10,dst=pagination.size
	fullname=items[].endAt,type=string,example=2023-01-30 00:00:00,dst=items[].endAt
	fullname=items[].title,type=string,example=新年好礼,dst=items[].title
	fullname=items,type=array,example=-,dst=items
	fullname=pagination,type=object,dst=pagination
	fullname=items[].valueObj,type=string,example=值对象,dst=items[].valueObj
	fullname=items[].type,type=string,example=image,dst=items[].type
	fullname=items[].beginAt,type=string,example=2023-01-12 00:00:00,dst=items[].beginAt
	fullname=items[].image,type=string,dst=items[].image
	fullname=items[].summary,type=string,example=下单有豪礼,dst=items[].summary
	fullname=items[].advertiserId,type=string,example=123,dst=items[].advertiserId
	fullname=code,type=string,example=-,dst=code
	fullname=items[].remark,type=string,example=营养早餐广告,dst=items[].remark
`
	l, err := jsonschemaline.ParseJsonschemaline(lineschema)
	if err != nil {
		panic(err)
	}
	jsonStr, err := l.Jsonschemaline2json()
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonStr)
}

func NewLineSchema() (l *jsonschemaline.Jsonschemaline) {
	var jsonStr = `
		[{
			"config":{
				"id":"1",
				"status":"2",
				"identify":"abcde",
				"merchantId":"123",
				"merchantName":"测试商户",
				"operateName":"彭政",
				"storeId":"1",
				"storeName":"门店名称",
				"ids":["1","2"],
				"array":[
					{"id":"2","name":"ok"}
					]
			}
		}]
	`
	lineschema, err := jsonschemaline.Json2lineSchema(jsonStr)
	if err != nil {
		panic(err)
	}
	return lineschema
}

func TestParse(t *testing.T) {
	line := `
	version=http://json-schema.org/draft-07/schema,id=input,direction=in
fullname=config.hsbRemark,dst=FhsbRemark,required
fullname=config.popUpWindow,dst=FpopUpWindow,format=number,required,enum=["0","1"]
fullname=config.xyRemark,dst=FxyRemark,required
fullname=config.status,dst=Fstatus,enum=["0","1"]
	`

	lineschema, err := jsonschemaline.ParseJsonschemaline(line)
	if err != nil {
		panic(err)
	}

	schema := new(jsonschemaline.Schema)
	schema.Raw2Schema(*lineschema)
	jsb, err := schema.MarshalJSON()
	if err != nil {
		panic(err)
	}
	jsonschemaStr := string(jsb)
	fmt.Println(jsonschemaStr)

}

func TestGjsonPath(t *testing.T) {
	line := `version=http://json-schema.org/draft-07/schema,id=output,direction=out
		fullname=items[].content,src=PaginateOut.#.content,required
		fullname=items[].createdAt,src=PaginateOut.#.created_at,required
		fullname=items[].deletedAt,src=PaginateOut.#.deleted_at,required
		fullname=items[].description,src=PaginateOut.#.description,required
		fullname=items[].icon,src=PaginateOut.#.icon,required
		fullname=items[].id,src=PaginateOut.#.id,required
		fullname=items[].key,src=PaginateOut.#.key,required
		fullname=items[].label,src=PaginateOut.#.label,required
		fullname=items[].thumb,src=PaginateOut.#.thumb,required
		fullname=items[].title,src=PaginateOut.#.title,required
		fullname=items[].updatedAt,src=PaginateOut.#.updated_at,required
		fullname=pageInfo.pageIndex,src=input.pageIndex,required
		fullname=pageInfo.pageSize,src=input.pageSize,required
		fullname=pageInfo.total,src=PaginateTotalOut,required`
	//{pageInfo:{pageIndex:input.pageIndex,pageSize:input.pageSize,total:PaginateTotalOut},items:{content:PaginateOut.#.content,createdAt:PaginateOut.#.created_at,deletedAt:PaginateOut.#.deleted_at}|@group}
	lineschema, err := jsonschemaline.ParseJsonschemaline(line)
	if err != nil {
		panic(err)
	}
	gjsonPath := lineschema.GjsonPath(nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(gjsonPath)
}
func TestGjsonPath2(t *testing.T) {
	// line := `version=http://json-schema.org/draft-07/schema,id=input,direction=in
	// fullname=pageIndex,dst=pageIndex,format=number,required
	// fullname=pageSize,dst=pageSize,format=number,required`
	//{pageInfo:{pageIndex:input.pageIndex,pageSize:input.pageSize,total:PaginateTotalOut},items:{content:PaginateOut.#.content,createdAt:PaginateOut.#.created_at,deletedAt:PaginateOut.#.deleted_at}|@group}
	// lineschema, err := jsonschemaline.ParseJsonschemaline(line)
	// if err != nil {
	// 	panic(err)
	// }
	// gjsonPath := lineschema.GjsonPath()
	// if err != nil {
	// 	panic(err)
	// }
	gjsonPath := "{pageIndex:input.pageIndex,pageSize:input.pageSize}"
	jsonStr := `{"input":{"pageIndex":"0","pageSize":"20"}}`
	out := gjson.Get(jsonStr, gjsonPath).String()
	fmt.Println(out)
}
func TestGjsonPathWithDefaultFormat(t *testing.T) {
	line := `version=http://json-schema.org/draft-07/schema,id=input,direction=in
	fullname=pageIndex,dst=pageIndex,format=number,required
	fullname=pageSize,dst=pageSize,format=number,required
	fullname=valid,dst=valid,format=bool,required
	`
	lineschema, err := jsonschemaline.ParseJsonschemaline(line)
	if err != nil {
		panic(err)
	}
	gjsonPath := lineschema.GjsonPathWithDefaultFormat()
	if err != nil {
		panic(err)
	}
	t.Run("string", func(t *testing.T) {
		jsonStr := `{"input":{"pageIndex":"0","pageSize":"20"}}`
		out := gjson.Get(jsonStr, gjsonPath).String()
		fmt.Println(out)
	})
	t.Run("nil", func(t *testing.T) {
		jsonStr := `{"input":{"pageSize":"20"}}`
		out := gjson.Get(jsonStr, gjsonPath).String()
		fmt.Println(out)
	})
	t.Run("true", func(t *testing.T) {
		jsonStr := `{"input":{"pageSize":"20","valid":1}}`
		out := gjson.Get(jsonStr, gjsonPath).String()
		fmt.Println(out)
	})

}

func TestJsonSchema(t *testing.T) {
	str := `{"Meta":{"id":"form","version":"http://json-schema.org/draft-07/schema#","direction":"in"},"Items":[{"comment":"广告标题","type":"string","required":"true","description":"广告标题","example":"新年豪礼","dst":"title","fullname":"title"},{"comment":"广告主","type":"string","required":"true","description":"广告主","example":"123","dst":"advertiserId","fullname":"advertiserId"},{"comment":"可以投放开始时间","type":"string","required":"true","description":"可以投放开始时间","example":"2023-01-12 00:00:00","dst":"beginAt","fullname":"beginAt"},{"comment":"投放结束时间","type":"string","required":"true","description":"投放结束时间","example":"2023-01-30 00:00:00","dst":"endAt","fullname":"endAt"},{"comment":"页索引,0开始","type":"string","required":"true","description":"页索引,0开始","default":"0","dst":"index","fullname":"index"},{"comment":"每页数量","type":"string","required":"true","description":"每页数量","default":"10","dst":"size","fullname":"size"},{"comment":"文件格式","type":"string","required":"true","description":"文件格式","default":"application/json","dst":"content-type","fullname":"content-type"},{"comment":"访问服务的备案id","type":"string","required":"true","description":"访问服务的备案id","dst":"appid","fullname":"appid"},{"comment":"签名,外网访问需开启签名","type":"string","required":"true","description":"签名,外网访问需开启签名","dst":"signature","fullname":"signature"}]}`
	lineSchema := jsonschemaline.Jsonschemaline{}
	json.Unmarshal([]byte(str), &lineSchema)
	b, err := lineSchema.JsonSchema()
	require.NoError(t, err)
	schema := string(b)
	expected := `{"$schema":"http://json-schema.org/draft-07/schema#","type":"object","required":["title","advertiserId","beginAt","endAt","index","size","content-type","appid","signature"],"properties":{"title":{"comment":"广告标题","type":"string","description":"广告标题","example":"新年豪礼"},"advertiserId":{"comment":"广告主","type":"string","description":"广告主","example":"123"},"beginAt":{"comment":"可以投放开始时间","type":"string","description":"可以投放开始时间","example":"2023-01-12 00:00:00"},"endAt":{"comment":"投放结束时间","type":"string","description":"投放结束时间","example":"2023-01-30 00:00:00"},"index":{"comment":"页索引,0开始","type":"string","description":"页索引,0开始","default":"0"},"size":{"comment":"每页数量","type":"string","description":"每页数量","default":"10"},"content-type":{"comment":"文件格式","type":"string","description":"文件格式","default":"application/json"},"appid":{"comment":"访问服务的备案id","type":"string","description":"访问服务的备案id"},"signature":{"comment":"签名,外网访问需开启签名","type":"string","description":"签名,外网访问需开启签名"}}}`
	ok := jsonpatch.Equal([]byte(expected), []byte(schema))
	assert.Equal(t, true, ok)
}

func TestToJsonSchemaKVS(t *testing.T) {
	t.Run("object", func(t *testing.T) {
		item := jsonschemaline.JsonschemalineItem{
			Fullname:    "._param.config.id",
			Description: "ID",
			Type:        "string",
			Format:      "number",
			Required:    true,
		}
		kvs, err := item.ToJsonSchemaKVS()
		require.NoError(t, err)
		schema, err := kvs.Json(false)
		require.NoError(t, err)
		expected := `{"$schema":"http://json-schema.org/draft-07/schema#","type":"object","properties":{"_param":{"type":"object","properties":{"config":{"type":"object","required":["id"],"properties":{"id":{"type":"string","format":"number","description":"ID"}}}}}}}`
		ok := jsonpatch.Equal([]byte(expected), []byte(schema))
		assert.Equal(t, true, ok)
	})

	t.Run("array", func(t *testing.T) {
		item := jsonschemaline.JsonschemalineItem{
			Fullname:    "[]._param.config.id",
			Description: "ID",
			Type:        "string",
			Format:      "number",
			Required:    true,
		}
		kvs, err := item.ToJsonSchemaKVS()
		require.NoError(t, err)
		schema, err := kvs.Json(false)
		require.NoError(t, err)
		expected := `{"$schema":"http://json-schema.org/draft-07/schema#","type":"array","items":{"type":"object","properties":{"_param":{"type":"object","properties":{"config":{"type":"object","required":["id"],"properties":{"id":{"type":"string","format":"number","description":"ID"}}}}}}}}`
		ok := jsonpatch.Equal([]byte(expected), []byte(schema))
		assert.Equal(t, true, ok)
	})
	t.Run("array_middle", func(t *testing.T) {
		item := jsonschemaline.JsonschemalineItem{
			Fullname:    "_param.config[].id",
			Description: "ID",
			Type:        "string",
			Format:      "number",
			Required:    true,
		}
		kvs, err := item.ToJsonSchemaKVS()
		require.NoError(t, err)
		schema, err := kvs.Json(false)
		require.NoError(t, err)
		expected := `{"$schema":"http://json-schema.org/draft-07/schema#","type":"object","properties":{"_param":{"type":"object","properties":{"config":{"type":"array","items":{"type":"object","required":["id"],"properties":{"id":{"type":"string","format":"number","description":"ID"}}}}}}}}`
		ok := jsonpatch.Equal([]byte(expected), []byte(schema))
		assert.Equal(t, true, ok)
	})
	t.Run("enumNames", func(t *testing.T) {
		item := jsonschemaline.JsonschemalineItem{
			Fullname:    "_param.config.type",
			Description: "类型",
			Type:        "string",
			Required:    true,
			Enum:        []string{"1", "2"},
			EnumNames:   []string{"类型1", "类型2"},
		}
		kvs, err := item.ToJsonSchemaKVS()
		require.NoError(t, err)
		schema, err := kvs.Json(false)
		require.NoError(t, err)
		fmt.Println(schema)
	})
	t.Run("enumNames_array", func(t *testing.T) {
		item := jsonschemaline.JsonschemalineItem{
			Fullname:    "_param.config.type[]",
			Description: "类型",
			Type:        "string",
			Required:    true,
			Enum:        []string{"1", "2"},
			EnumNames:   []string{"类型1", "类型2"},
		}
		kvs, err := item.ToJsonSchemaKVS()
		require.NoError(t, err)
		schema, err := kvs.Json(false)
		require.NoError(t, err)
		fmt.Println(schema)
	})
}

func TestToSturct(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		str := `{"Meta":{"id":"form","version":"http://json-schema.org/draft-07/schema#","direction":"in"},"Items":[{"comment":"广告标题","type":"string","required":"true","description":"广告标题","example":"新年豪礼","dst":"title","fullname":"title"},{"comment":"广告主","type":"string","required":"true","description":"广告主","example":"123","dst":"advertiserId","fullname":"advertiserId"},{"comment":"可以投放开始时间","type":"string","required":"true","description":"可以投放开始时间","example":"2023-01-12 00:00:00","dst":"beginAt","fullname":"beginAt"},{"comment":"投放结束时间","type":"string","required":"true","description":"投放结束时间","example":"2023-01-30 00:00:00","dst":"endAt","fullname":"endAt"},{"comment":"页索引,0开始","type":"string","required":"true","description":"页索引,0开始","default":"0","format":"int","dst":"index","fullname":"index"},{"comment":"每页数量","type":"string","required":"true","description":"每页数量","default":"10","dst":"size","fullname":"size"},{"comment":"文件格式","type":"string","required":"true","description":"文件格式","default":"application/json","dst":"content-type","fullname":"content-type"},{"comment":"访问服务的备案id","type":"string","required":"true","description":"访问服务的备案id","dst":"appid","fullname":"appid"},{"comment":"签名,外网访问需开启签名","type":"string","required":"true","description":"签名,外网访问需开启签名","dst":"signature","fullname":"signature"}]}`
		lineSchema := jsonschemaline.Jsonschemaline{}
		json.Unmarshal([]byte(str), &lineSchema)
		structs := lineSchema.ToSturct()
		fmt.Println(structs.Json())
	})

	t.Run("complex", func(t *testing.T) {
		str := `
		version=http://json-schema.org/draft-07/schema#,direction=out,id=out
        fullname=code,src=code,description=业务状态码,comment=业务状态码,example=0
        fullname=message,src=message,description=业务提示,comment=业务提示,example=ok
        fullname=items,src=items,type=array,description=数组,comment=数组,example=-
        fullname=items[].id,src=items[].id,description=主键,comment=主键,example=0
        fullname=items[].title,src=items[].title,description=广告标题,comment=广告标题,example=新年豪礼
        fullname=items[].advertiserId,src=items[].advertiserId,description=广告主,comment=广告主,example=123
        fullname=items[].summary,src=items[].summary,description=广告素材-文字描述,comment=广告素材-文字描述,example=下单有豪礼
        fullname=items[].image,src=items[].image,description=广告素材-图片地址,comment=广告素材-图片地址
        fullname=items[].link,src=items[].link,description=连接地址,comment=连接地址
        fullname=items[].type,src=items[].type,description=广告素材(类型),text-文字,image-图片,vido-视频,comment=广告素材(类型),text-文字,image-图片,vido-视频,example=image
        fullname=items[].beginAt,src=items[].beginAt,description=投放开始时间,comment=投放开始时间,example=2023-01-12 00:00:00
        fullname=items[].endAt,src=items[].endAt,description=投放结束时间,comment=投放结束时间,example=2023-01-30 00:00:00
        fullname=items[].remark,src=items[].remark,description=备注,comment=备注,example=营养早餐广告
        fullname=items[].valueObj,src=items[].valueObj,description=json扩展,广告的值属性对象,comment=json扩展,广告的值属性对象,example={"tag":"index"}
        fullname=pagination,src=pagination,type=object,description=对象,comment=对象
        fullname=pagination.index,src=pagination.index,description=页索引,0开始,comment=页索引,0开始,example=0
        fullname=pagination.size,src=pagination.size,description=每页数量,comment=每页数量,example=10
        fullname=pagination.total,src=pagination.total,description=总数,comment=总数,example=60
`
		lineSchema, err := jsonschemaline.ParseJsonschemaline(str)
		require.NoError(t, err)
		structs := lineSchema.ToSturct()
		b, err := json.Marshal(structs)
		require.NoError(t, err)
		fmt.Println(string(b))
	})

}

func TestExampleSjsn(t *testing.T) {
	lineschema := "version=http://json-schema.org/draft-07/schema#,direction=in,id=example\nfullname=[[Table . `obj.#(fullname%\"items[]*\")#,dst=[[Table . `obj.#(fullname%\"items[]*\")#,type=#(fullname!%\"*[].id\")#` \",format={#.fullname.@basePath},default=是,example={#.description}"
	schema, err := jsonschemaline.ParseJsonschemaline(lineschema)
	require.NoError(t, err)
	example, err := schema.JsonExample()
	require.NoError(t, err)
	fmt.Println(example)
}
func TestAddNameprefix(t *testing.T) {
	str := `
		version=http://json-schema.org/draft-07/schema#,direction=out,id=out
        fullname=code,src=code,description=业务状态码,comment=业务状态码,example=0
        fullname=message,src=message,description=业务提示,comment=业务提示,example=ok
        fullname=items,src=items,type=array,description=数组,comment=数组,example=-
        fullname=items[].id,src=items[].id,description=主键,comment=主键,example=0
        fullname=items[].title,src=items[].title,description=广告标题,comment=广告标题,example=新年豪礼
        fullname=items[].advertiserId,src=items[].advertiserId,description=广告主,comment=广告主,example=123
        fullname=items[].summary,src=items[].summary,description=广告素材-文字描述,comment=广告素材-文字描述,example=下单有豪礼
        fullname=items[].image,src=items[].image,description=广告素材-图片地址,comment=广告素材-图片地址
        fullname=items[].link,src=items[].link,description=连接地址,comment=连接地址
        fullname=items[].type,src=items[].type,description=广告素材(类型),text-文字,image-图片,vido-视频,comment=广告素材(类型),text-文字,image-图片,vido-视频,example=image
        fullname=items[].beginAt,src=items[].beginAt,description=投放开始时间,comment=投放开始时间,example=2023-01-12 00:00:00
        fullname=items[].endAt,src=items[].endAt,description=投放结束时间,comment=投放结束时间,example=2023-01-30 00:00:00
        fullname=items[].remark,src=items[].remark,description=备注,comment=备注,example=营养早餐广告
        fullname=items[].valueObj,src=items[].valueObj,description=json扩展,广告的值属性对象,comment=json扩展,广告的值属性对象,example={"tag":"index"}
        fullname=pagination,src=pagination,type=object,description=对象,comment=对象
        fullname=pagination.index,src=pagination.index,description=页索引,0开始,comment=页索引,0开始,example=0
        fullname=pagination.size,src=pagination.size,description=每页数量,comment=每页数量,example=10
        fullname=pagination.total,src=pagination.total,description=总数,comment=总数,example=60
`
	lineSchema, err := jsonschemaline.ParseJsonschemaline(str)
	require.NoError(t, err)
	structs := lineSchema.ToSturct()
	nameprefix := "name_space"
	newStructs := structs.Copy()
	newStructs.AddNameprefix(nameprefix)
	fmt.Println(newStructs)
}
