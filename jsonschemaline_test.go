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

func TestParseDefaultJson(t *testing.T) {
	var schemalineIn3 = `
    version=http://json-schema.org/draft-07/schema#,direction=in,id=mainIn
	fullname=pageSize,format=number,required,dst=Limit,default=20
fullname=pageIndex,format=number,required,dst=pageIndex,default=1
 `
	schemaline, err := jsonschemaline.ParseJsonschemaline(schemalineIn3)
	if err != nil {
		panic(err)
	}
	defaultJson, err := schemaline.DefaultJson()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v", defaultJson)
}

func TestJsonSchemalineString(t *testing.T) {
	lineschema, err := jsonschemaline.ParseJsonschemaline(schemalineOut)
	if err != nil {
		panic(err)
	}

	fmt.Println(lineschema.String())

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
				"storeName":"门店名称",
				"more":[
					{"key":"1","value":"张三"},
					{"key":"2","value":"李四"}
				]
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
	jsonStr, err := l.JsonExample()
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
					{"id":"2","name":"ok"},
					{"id":"3","name":"ok"}
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
	gjsonPath := lineschema.GjsonPath(false, nil)
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

func TestGjsonPathWithDefaultFormatOutput(t *testing.T) {
	inputLineschema := `version=http://json-schema.org/draft-07/schema,id=output,direction=out
	fullname=pageIndex,src=pageIndex,format=number,required
	fullname=pageSize,src=pageSize,format=number,required
	fullname=valid,src=valid,format=bool,required
	fullname=items[].id,src=items.#.id,format=int,required
	fullname=items[].title,src=items.#.title,required
	fullname=items[].status,src=items.#.status,format=bool,required
	`
	lineschema, err := jsonschemaline.ParseJsonschemaline(inputLineschema)
	if err != nil {
		panic(err)
	}
	t.Run("string", func(t *testing.T) {
		jsonStr := `{"pageIndex":"0","pageSize":"20","items":[{"id":"1","title":"标题1","status":true},{"id":"2","title":"标题2","status":false}]}`
		gjsonPath := lineschema.GjsonPathWithDefaultFormat(true)
		if err != nil {
			panic(err)
		}
		out := gjson.Get(jsonStr, gjsonPath).String()
		excepted := `{"pageIndex":"0","pageSize":"20","items":[{"status":"true","id":"1","title":"标题1"},{"status":"false","id":"2","title":"标题2"}]}`
		assert.JSONEq(t, excepted, out)
	})
}
func TestGjsonPathWithDefaultFormatInput(t *testing.T) {
	inputLineschema := `version=http://json-schema.org/draft-07/schema,id=input,direction=in
	fullname=pageIndex,dst=pageIndex,format=number,required
	fullname=pageSize,dst=pageSize,format=number,required
	fullname=valid,dst=valid,format=bool,required
	fullname=items[].id,dst=items.#.id,format=int,required
	fullname=items[].title,dst=items.#.title,required
	`
	lineschema, err := jsonschemaline.ParseJsonschemaline(inputLineschema)
	if err != nil {
		panic(err)
	}

	t.Run("string", func(t *testing.T) {
		jsonStr := `{"pageIndex":"0","pageSize":"20","items":[{"id":"1","title":"标题1"},{"id":"2","title":"标题2"}]}`
		gjsonPath := lineschema.GjsonPathWithDefaultFormat(false)
		if err != nil {
			panic(err)
		}
		out := gjson.Get(jsonStr, gjsonPath).String()
		fmt.Println(out)
	})
	t.Run("nil", func(t *testing.T) {
		jsonStr := `{"pageSize":"20"}`
		gjsonPath := lineschema.GjsonPathWithDefaultFormat(false)
		if err != nil {
			panic(err)
		}
		out := gjson.Get(jsonStr, gjsonPath).String()
		fmt.Println(out)
	})
	t.Run("true", func(t *testing.T) {
		jsonStr := `{"pageSize":"20","valid":1}`
		gjsonPath := lineschema.GjsonPathWithDefaultFormat(false)
		if err != nil {
			panic(err)
		}
		out := gjson.Get(jsonStr, gjsonPath).String()
		fmt.Println(out)
	})

	t.Run("true-ingnoreID", func(t *testing.T) {
		jsonStr := `{"pageSize":"20","valid":1}`
		gjsonPath := lineschema.GjsonPathWithDefaultFormat(true)
		if err != nil {
			panic(err)
		}
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
			Enum:        `[1,2]`,
			EnumNames:   `["类型1", "类型2"]`,
		}
		kvs, err := item.ToJsonSchemaKVS()
		require.NoError(t, err)
		schema, err := kvs.Json(true)
		require.NoError(t, err)
		fmt.Println(schema)
	})
	t.Run("enumNames_array", func(t *testing.T) {
		item := jsonschemaline.JsonschemalineItem{
			Fullname:    "_param.config.type[]",
			Description: "类型",
			Type:        "string",
			Required:    true,
			Enum:        `[1,2]`,
			EnumNames:   `["类型1", "类型2"]`,
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

func TestGetJsonSchemaSchema(t *testing.T) {
	schema := jsonschemaline.GetJsonSchemaSchema()
	fmt.Println(schema)
}

func TestFormatPath(t *testing.T) {
	data := `{"code":"","message":"","data":[{"fullname":"fullname","type":"string","title":"名称/全称","description":"名称/全称","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"type","type":"string","title":"类型","description":"类型","format":"","more":{"enum":["int","float","string"],"enumNames":["整型","浮点型","字符串"],"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"enum","type":"string","title":"枚举值","description":"枚举值","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"enumNames","type":"string","title":"枚举值标题","description":"枚举值标题","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"comment","type":"string","title":"备注","description":"备注","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"const","type":"string","title":"常量","description":"常量","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"multipleOf","type":"string","title":"多值","description":"多值","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"maximum","type":"string","title":"最大值","description":"最大值","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"exclusiveMaximum","type":"string","title":"是否包含最大值","description":"是否包含最大值","format":"bool","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"minimum","type":"string","title":"最小值","description":"最小值","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"exclusiveMinimum","type":"string","title":"是否包含最小值","description":"是否包含最小值","format":"bool","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"maxLength","type":"string","title":"最大长度","description":"最大长度","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"minLength","type":"string","title":"最小长度","description":"最小长度","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"pattern","type":"string","title":"匹配格式","description":"匹配格式","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"maxItems","type":"string","title":"最大项数","description":"最大项数","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"minItems","type":"string","title":"最小项数","description":"最小项数","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"uniqueItems","type":"string","title":"数组唯一","description":"数组唯一","format":"bool","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"maxContains","type":"string","title":"符合Contains规则最大数量","description":"符合Contains规则最大数量","format":"uint","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"minContains","type":"string","title":"符合Contains规则最小数量","description":"符合Contains规则最小数量","format":"uint","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"maxProperties","type":"string","title":"对象最多属性个数","description":"对象最多属性个数","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"minProperties","type":"string","title":"对象最少属性个数","description":"对象最少属性个数","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"required","type":"string","title":"是否必须","description":"是否必须","format":"bool","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"format","type":"string","title":"类型格式","description":"类型格式","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"contentEncoding","type":"string","title":"内容编码","description":"内容编码","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"contentMediaType","type":"string","title":"内容格式","description":"内容格式","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"title","type":"string","title":"标题","description":"标题","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"description","type":"string","title":"描述","description":"描述","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"default","type":"string","title":"默认值","description":"默认值","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"deprecated","type":"string","title":"是否弃用","description":"是否弃用","format":"bool","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"readOnly","type":"string","title":"只读","description":"只读","format":"bool","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"writeOnly","type":"string","title":"只写","description":"只写","format":"bool","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"example","type":"string","title":"案例","description":"案例","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"examples","type":"string","title":"案例集合","description":"案例集合","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"allowEmptyValue","type":"string","title":"是否可以为空","description":"是否可以为空","format":"bool","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}}]}`
	path := `{message:message.@tostring,data:{format:data.#.format.@tostring,more:{exclusiveMaximum:data.#.more.exclusiveMaximum.@tostring,minProperties:data.#.more.minProperties.@tostring,multipleOf:data.#.more.multipleOf.@tostring,maxLength:data.#.more.maxLength.@tostring,minContains:data.#.more.minContains.@tostring,minimum:data.#.more.minimum.@tostring,allowEmptyValue:data.#.more.allowEmptyValue.@tostring,pattern:data.#.more.pattern.@tostring,maxItems:data.#.more.maxItems.@tostring,default:data.#.more.default.@tostring,writeOnly:data.#.more.writeOnly.@tostring,const:data.#.more.const.@tostring,contentEncoding:data.#.more.contentEncoding.@tostring,readOnly:data.#.more.readOnly.@tostring,example:data.#.more.example.@tostring,exclusiveMinimum:data.#.more.exclusiveMinimum.@tostring,required:data.#.more.required.@tostring,contentMediaType:data.#.more.contentMediaType.@tostring,deprecated:data.#.more.deprecated.@tostring,comment:data.#.more.comment.@tostring,minLength:data.#.more.minLength.@tostring,minItems:data.#.more.minItems.@tostring,fullname:data.#.more.fullname.@tostring,examples:data.#.more.examples.@tostring,maximum:data.#.more.maximum.@tostring,uniqueItems:data.#.more.uniqueItems.@tostring,maxContains:data.#.more.maxContains.@tostring,maxProperties:data.#.more.maxProperties.@tostring}|@group,fullname:data.#.fullname.@tostring,type:data.#.type.@tostring,title:data.#.title.@tostring,description:data.#.description.@tostring}|@group,code:code.@tostring}`
	out := gjson.Get(data, path).String()
	fmt.Println(out)
}

func TestGjsonPathWithDefaultFormat(t *testing.T) {
	outputLineschema := `version=http://json-schema.org/draft-07/schema#,direction=out,id=out
	fullname=code,src=code,title=业务状态码,comment=业务状态码,example=0
	fullname=message,src=message,title=业务提示,comment=业务提示,example=ok
	fullname=data[].fullname,src=data.#.fullname,title=名称,comment=名称
	fullname=data[].type,src=data.#.type,title=类型"
	fullname=data[].title,src=data.#.title,title=标题,comment=标题
	fullname=data[].description,src=data.#.description,title=描述,comment=描述
	fullname=data[].format,src=data.#.format,title=类型格式,comment=类型格式
	fullname=data[].more.enum[],src=data.#.more.enum.#,type=string,title=枚举值,comment=枚举值
	fullname=data[].more.enumNames[],src=data.#.more.enumNames.#,type=string,title=枚举值标题,comment=枚举值标题
	fullname=data[].more.comment,src=data.#.more.comment,title=备注,comment=备注
	fullname=data[].more.const,src=data.#.more.const,title=常量,comment=常量
	fullname=data[].more.multipleOf,src=data.#.more.multipleOf,format=int,title=多值,comment=多值
	fullname=data[].more.maximum,src=data.#.more.maximum,format=int,title=最大值,comment=最大值
	fullname=data[].more.exclusiveMaximum,src=data.#.more.exclusiveMaximum,format=bool,title=是否包含最大值,comment=是否包含最大值
	fullname=data[].more.minimum,src=data.#.more.minimum,format=int,title=最小值,comment=最小值
	fullname=data[].more.exclusiveMinimum,src=data.#.more.exclusiveMinimum,format=bool,title=是否包含最小值,comment=是否包含最小值
	fullname=data[].more.maxLength,src=data.#.more.maxLength,format=int,title=最大长度,comment=最大长度
	fullname=data[].more.minLength,src=data.#.more.minLength,format=int,title=最小长度,comment=最小长度
	fullname=data[].more.pattern,src=data.#.more.pattern,title=匹配格式,comment=匹配格式
	fullname=data[].more.maxItems,src=data.#.more.maxItems,format=int,title=最大项数,comment=最大项数
	fullname=data[].more.minItems,src=data.#.more.minItems,format=int,title=最小项数,comment=最小项数
	fullname=data[].more.uniqueItems,src=data.#.more.uniqueItems,format=bool,title=数组唯一,comment=数组唯一
	fullname=data[].more.maxContains,src=data.#.more.maxContains,format=uint,title=符合Contains规则最大数量,comment=符合Contains规则最大数量
	fullname=data[].more.minContains,src=data.#.more.minContains,format=uint,title=符合Contains规则最小数量,comment=符合Contains规则最小数量
	fullname=data[].more.maxProperties,src=data.#.more.maxProperties,format=int,title=对象最多属性个数,comment=对象最多属性个数
	fullname=data[].more.minProperties,src=data.#.more.minProperties,format=int,title=对象最少属性个数,comment=对象最少属性个数
	fullname=data[].more.required,src=data.#.more.required,format=bool,title=是否必须,comment=是否必须
	fullname=data[].more.contentEncoding,src=data.#.more.contentEncoding,title=内容编码,comment=内容编码
	fullname=data[].more.contentMediaType,src=data.#.more.contentMediaType,title=内容格式,comment=内容格式
	fullname=data[].more.default,src=data.#.more.default,title=默认值,comment=默认值
	fullname=data[].more.deprecated,src=data.#.more.deprecated,format=bool,title=是否弃用,comment=是否弃用
	fullname=data[].more.readOnly,src=data.#.more.readOnly,format=bool,title=只读,comment=只读
	fullname=data[].more.writeOnly,src=data.#.more.writeOnly,format=bool,title=只写,comment=只写
	fullname=data[].more.example,src=data.#.more.example,title=案例,comment=案例
	fullname=data[].more.examples,src=data.#.more.examples,title=案例集合,comment=案例集合
	fullname=data[].more.fullname,src=data.#.more.fullname,title=名称/全称,comment=名称/全称
	fullname=data[].more.allowEmptyValue,src=data.#.more.allowEmptyValue,format=bool,title=是否可以为空,comment=是否可以为空`
	data := `{"code":"","message":"","data":[{"fullname":"fullname","type":"string","title":"名称/全称","description":"名称/全称","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"type","type":"string","title":"类型","description":"类型","format":"","more":{"enum":["int","float","string"],"enumNames":["整型","浮点型","字符串"],"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"enum","type":"string","title":"枚举值","description":"枚举值","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"enumNames","type":"string","title":"枚举值标题","description":"枚举值标题","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"comment","type":"string","title":"备注","description":"备注","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"const","type":"string","title":"常量","description":"常量","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"multipleOf","type":"string","title":"多值","description":"多值","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"maximum","type":"string","title":"最大值","description":"最大值","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"exclusiveMaximum","type":"string","title":"是否包含最大值","description":"是否包含最大值","format":"bool","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"minimum","type":"string","title":"最小值","description":"最小值","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"exclusiveMinimum","type":"string","title":"是否包含最小值","description":"是否包含最小值","format":"bool","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"maxLength","type":"string","title":"最大长度","description":"最大长度","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"minLength","type":"string","title":"最小长度","description":"最小长度","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"pattern","type":"string","title":"匹配格式","description":"匹配格式","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"maxItems","type":"string","title":"最大项数","description":"最大项数","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"minItems","type":"string","title":"最小项数","description":"最小项数","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"uniqueItems","type":"string","title":"数组唯一","description":"数组唯一","format":"bool","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"maxContains","type":"string","title":"符合Contains规则最大数量","description":"符合Contains规则最大数量","format":"uint","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"minContains","type":"string","title":"符合Contains规则最小数量","description":"符合Contains规则最小数量","format":"uint","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"maxProperties","type":"string","title":"对象最多属性个数","description":"对象最多属性个数","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"minProperties","type":"string","title":"对象最少属性个数","description":"对象最少属性个数","format":"int","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"required","type":"string","title":"是否必须","description":"是否必须","format":"bool","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"format","type":"string","title":"类型格式","description":"类型格式","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"contentEncoding","type":"string","title":"内容编码","description":"内容编码","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"contentMediaType","type":"string","title":"内容格式","description":"内容格式","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"title","type":"string","title":"标题","description":"标题","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"description","type":"string","title":"描述","description":"描述","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"default","type":"string","title":"默认值","description":"默认值","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"deprecated","type":"string","title":"是否弃用","description":"是否弃用","format":"bool","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"readOnly","type":"string","title":"只读","description":"只读","format":"bool","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"writeOnly","type":"string","title":"只写","description":"只写","format":"bool","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"example","type":"string","title":"案例","description":"案例","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"examples","type":"string","title":"案例集合","description":"案例集合","format":"","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}},{"fullname":"allowEmptyValue","type":"string","title":"是否可以为空","description":"是否可以为空","format":"bool","more":{"enum":null,"enumNames":null,"comment":"","const":"","multipleOf":0,"maximum":0,"exclusiveMaximum":false,"minimum":0,"exclusiveMinimum":false,"maxLength":0,"minLength":0,"pattern":"","maxItems":0,"minItems":0,"uniqueItems":false,"maxContains":0,"minContains":0,"maxProperties":0,"minProperties":0,"required":false,"contentEncoding":"","contentMediaType":"","default":"","deprecated":false,"readOnly":false,"writeOnly":false,"example":"","examples":"","fullname":"","allowEmptyValue":false}}]}`
	jsonline, err := jsonschemaline.ParseJsonschemaline(outputLineschema)
	require.NoError(t, err)
	gpath := jsonline.GjsonPathWithDefaultFormat(true)
	newData := gjson.Get(data, gpath).String()
	fmt.Println(newData)

}

func TestGjsonPathWithDefaultFormatOutput2(t *testing.T) {
	outputLineschema := `version=http://json-schema.org/draft-07/schema#,direction=out,id=out
fullname=reportUrl,src=reportUrl,title=质检报告URL,comment=质检报告URL
fullname=classId,src=classId,format=int,title=分类ID,comment=分类ID
fullname=className,src=className,title=分类名称,comment=分类名称
fullname=productId,src=productId,format=int,title=产品ID,comment=产品ID
fullname=productName,src=productName,title=产品名称,comment=产品名称
fullname=brandId,src=brandId,format=int,title=品牌id,comment=品牌id
fullname=brandName,src=brandName,title=品牌名称,comment=品牌名称
fullname=detectType,src=detectType,title=检测类型line插线检imei扫描IMEI号manual手动检,comment=检测类型line插线检imei扫描IMEI号manual手动检
fullname=imeiResult,src=imeiResult,type=object,title=imei检测结果,comment=imei检测结果
fullname=imeiResult.imei,src=imeiResult.imei,title=imei号,comment=imei号
fullname=imeiResult.items,src=imeiResult.items,type=array,title=imei选项,comment=imei选项
fullname=imeiResult.items[].qid,src=imeiResult.items.#.qid,title=问题ID,comment=问题ID
fullname=imeiResult.items[].qname,src=imeiResult.items.#.qname,title=问题名称,comment=问题名称
fullname=imeiResult.items[].aname,src=imeiResult.items.#.aname,title=答案名称,comment=答案名称
fullname=SelectedAnswers,src=SelectedAnswers,type=array,title=答案,comment=答案
fullname=SelectedAnswers,src=SelectedAnswers,type=array,title=选中的答案集合,comment=选中的答案集合
fullname=SelectedAnswers[].selectorType,src=SelectedAnswers.#.selectorType,title=答题者类型,comment=答题者类型
fullname=SelectedAnswers[].selectorId,src=SelectedAnswers.#.selectorId,title=答题者id,comment=答题者id
fullname=SelectedAnswers[].aId,src=SelectedAnswers.#.aId,title=选中的答案ID,comment=选中的答案ID
fullname=SelectedAnswers[].aname,src=SelectedAnswers.#.aname,title=答案项名称,comment=答案项名称
fullname=SelectedAnswers[].qname,src=SelectedAnswers.#.qname,title=问题项名称,comment=问题项名称
fullname=SelectedAnswers[].qId,src=SelectedAnswers.#.qId,title=选中的问题ID,comment=选中的问题ID
fullname=SelectedAnswers[].stepName,src=SelectedAnswers.#.stepName,title=所属步骤名称,comment=所属步骤名称
fullname=SelectedAnswers[].isBad,src=SelectedAnswers.#.isBad,format=int,title=是否是缺陷项（1-是,0-不是）,comment=是否是缺陷项（1-是,0-不是）`

	outputLineSchema, err := jsonschemaline.ParseJsonschemaline(outputLineschema)
	require.NoError(t, err)
	outputFormatGjsonPath := outputLineSchema.GjsonPathWithDefaultFormat(true)
	fmt.Println(outputFormatGjsonPath)

}

func TestArray(t *testing.T) {
	outputLineschema := `version=http://json-schema.org/draft-07/schema#,direction=out,id=out
fullname=times[].time,src=times.#.time,title=时间,comment=时间
fullname=select,src=select,type=array,format=string,title=17项选项ID,comment=17项选项ID`
	outputLineSchema, err := jsonschemaline.ParseJsonschemaline(outputLineschema)
	require.NoError(t, err)
	outputFormatGjsonPath := outputLineSchema.GjsonPathWithDefaultFormat(true)
	fmt.Println(outputFormatGjsonPath)

}

func TestMergeLineschema(t *testing.T) {
	t.Run("in", func(t *testing.T) {
		first := `
		version=http://json-schema.org/draft-07/schema,id=input,direction=in
	fullname=config.id,dst=id,required
	fullname=config.keyConst,dst=keyConst,required
	fullname=config.label,dst=label,required
	fullname=config.title,dst=title,required
	fullname=pageInfo.pageIndex,dst=pageInfo.pageIndex,required
	fullname=pageInfo.pageSize,dst=pageInfo.pageSize,required
		
		`
		second := `
		version=http://json-schema.org/draft-07/schema,id=input,direction=in
		fullname=Fid,type=int,dst=Fid,required
		fullname=Fkey_const,dst=Fkey_const,required
		fullname=Flabel,dst=Flabel,required
		fullname=Ftitle,dst=title,required
		`
		merged, err := jsonschemaline.MergeLineschema(first, second, 0.7)
		require.NoError(t, err)
		fmt.Println(merged)
	})

	t.Run("out", func(t *testing.T) {
		first := `
	version=http://json-schema.org/draft-07/schema,id=output,direction=out
fullname=items[].content,src=items.#.content,required
fullname=items[].createdAt,src=items.#.created_at,required
fullname=items[].deletedAt,src=items.#.deleted_at,required
fullname=items[].description,src=items.#.description,required
fullname=items[].icon,src=items.#.icon,required
fullname=items[].id,src=items.#.id,required
fullname=items[].key,src=items.#.key,required
fullname=items[].label,src=items.#.label,required
fullname=items[].thumb,src=items.#.thumb,required
fullname=items[].title,src=items.#.title,required
fullname=items[].updatedAt,src=items.#.updated_at,required
fullname=pageInfo.pageIndex,src=input.pageInd#.ex,required
fullname=pageInfo.pageSize,src=input.pageSize,required
fullname=pageInfo.total,src=PaginateTotalOut,required
	`
		second := `
	version=http://json-schema.org/draft-07/schema,id=output,direction=out
fullname=PaginateOut[].Fcontent,src=PaginateOut.#.Fcontent,required
fullname=PaginateOut[].FcreatedAt,src=PaginateOut.#.Fcreated_at,required
fullname=PaginateOut[].FdeletedAt,src=PaginateOut.#.Fdeleted_at,required
fullname=PaginateOut[].Fdescription,src=PaginateOut.#.Fdescription,required
fullname=PaginateOut[].Ficon,src=PaginateOut.#.Ficon,required
fullname=PaginateOut[].Fid,src=PaginateOut.#.Fid,required
fullname=PaginateOut[].Fkey,src=PaginateOut.#.Fkey,required
fullname=PaginateOut[].Flabel,src=PaginateOut.#.Flabel,required
fullname=PaginateOut[].Fthumb,src=PaginateOut.#.Fthumb,required
fullname=PaginateOut[].Ftitle,src=PaginateOut.#.title,required
fullname=PaginateOut[].FupdatedAt,src=PaginateOut.#.updated_at,required
	
	`

		merged, err := jsonschemaline.MergeLineschema(first, second, 0.7)
		require.NoError(t, err)
		fmt.Println(merged)
	})

}
