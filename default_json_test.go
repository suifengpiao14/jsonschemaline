package jsonschemaline_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/jsonschemaline"
)

var schemalineIn3 = `
    version=http://json-schema.org/draft-07/schema#,id=mainIn
	fullname=pageSize,format=number,required,dst=Limit
fullname=pageIndex,format=number,required,dst={{setValue . "Offset" (mul  (getValue .  "mainIn.pageIndex")   (getValue . "mainIn.pageSize"))}}
 `

func TestParseDefaultJson(t *testing.T) {

	schemaline, err := jsonschemaline.ParseJsonschemaline(schemalineIn3)
	if err != nil {
		panic(err)
	}
	defaultJson, err := jsonschemaline.ParseDefaultJson(*schemaline)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v", defaultJson)
}

func TestJsonMerge(t *testing.T) {
	defaultJson := `{"pageSize":"20","remark":"hello world"}`
	specialJson := `{"pageIndex":"0","pageSize":"10"}`

	merge, err := jsonschemaline.JsonMerge(defaultJson, specialJson)
	if err != nil {
		panic(err)
	}
	fmt.Println(merge)
}
