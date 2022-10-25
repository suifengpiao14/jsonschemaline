package jsonschemaline_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/jsonschemaline"
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
fullname=items[].amcUserName,src=hjxAmcRelationPaginateOut.#.Freal_name,required
fullname=items[].amcUserId,src=hjxAmcRelationPaginateOut.#.Famc_id,required
fullname=items[].amcUserEmail,src=hjxAmcRelationPaginateOut.#.Famc_name,required
fullname=items[].xyxzUserId,src=hjxAmcRelationPaginateOut.#.Fuser_id,required
fullname=pageInfo.pageIndex,src=input.pageIndex,required
fullname=pageInfo.pageSize,src=input.pageSize,required
fullname=pageInfo.total,src=hjxAmcRelationTotalOut,required`
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
