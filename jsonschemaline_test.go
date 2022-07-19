package jsonschemaline_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/jsonschemaline"
)

var jsonStr = `{"$schema":"http://json-schema.org/draft-07/schema#","$id":"main","type":"object","properties":{"config":{"type":"object","properties":{"hsbRemark":{"type":"string"},"id":{"type":"string","format":"number"},"popUpWindow":{"type":"string","enum":["0","1"]},"xyRemark":{"type":"string"},"status":{"type":"string","enum":["0","1"]}},"required":["hsbRemark","id","popUpWindow","xyRemark"]}},"required":["config"]}`

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
