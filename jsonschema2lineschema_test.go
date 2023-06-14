package jsonschemaline_test

import (
	"testing"

	"github.com/suifengpiao14/jsonschemaline"
)

func TestJsonSchema2LineSchema(t *testing.T) {
	var jsonschmaStr = `{"$schema":"http://json-schema.org/draft-07/schema#","$id":"execAPIInquiryScreenIdentifyUpdate","type":"object","required":["config"],"properties":{"config":{"properties":{"id":{"type":"string","format":"number"},"status":{"type":"string","enum":["1","2"]},"identify":{"type":"string","minLength":1},"merchantId":{"type":"string","format":"number"},"merchantName":{"type":"string","minLength":1},"operateName":{"type":"string","minLength":1},"storeId":{"type":"string","format":"number"},"storeName":{"type":"string","minLength":1}},"type":"object","required":["id","status","identify","merchantId","merchantName","operateName","storeId","storeName"]}}}`

	lineschemaStr, err := jsonschemaline.JsonSchema2LineSchema(jsonschmaStr)
	if err != nil {
		panic(err)
	}

	_ = lineschemaStr
	//fmt.Println(lineschemaStr)

}
