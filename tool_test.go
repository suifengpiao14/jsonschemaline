package jsonschemaline_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/jsonschemaline"
)

func TestJsonMerge(t *testing.T) {
	defaultJson := `{"pageSize":"20","remark":"hello world"}`
	specialJson := `{"pageIndex":"0","pageSize":""}`

	merge, err := jsonschemaline.JsonMergeIgnoreEmptyString(defaultJson, specialJson)
	if err != nil {
		panic(err)
	}
	fmt.Println(merge)
}
