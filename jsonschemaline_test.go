package jsonschemaline_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/jsonschemaline"
)

func TestJsonSchemalineString(t *testing.T) {
	lineschema, err := jsonschemaline.ParseJsonschemaline(schemalineOut)
	if err != nil {
		panic(err)
	}

	fmt.Println(lineschema.String())

}
