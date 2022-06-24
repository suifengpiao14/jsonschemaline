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

	defaultJsons, err := jsonschemaline.ParseDefaultJson(schemalineIn3)
	if err != nil {
		panic(err)
	}
	for _, defaultJson := range defaultJsons {
		fmt.Printf("%#v", defaultJson)
	}
}
