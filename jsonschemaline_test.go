package jsonschemaline_test

import (
	"fmt"
	"testing"

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
func TestJsonschemaline2json(t *testing.T) {

	l := NewLineSchema()
	jsonStr, err := l.Jsonschemaline2json()
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonStr)
}
func TestJsonschemaline2json2(t *testing.T) {
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
	jsonStr, err := l.Jsonschemaline2json()
	if err != nil {
		panic(err)
	}
	fmt.Println(jsonStr)
}

func NewLineSchema() (l *jsonschemaline.Jsonschemaline) {
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
				"array":[
					{"id":"2","name":"ok"}
					]
			}
		}
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
