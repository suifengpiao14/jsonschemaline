package jsonschemaline_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/jsonschemaline"
)

var schemalineOut = `
version=http://json-schema.org/draft-07/schema#,id=hello_world
fullname=items[].id,src=PaginateOut.#.Fid,required
fullname=items[].openId,src=PaginateOut.#.Fopen_id,required
fullname=items[].type,src=PaginateOut.#.Fopen_id_type,required
fullname=items[].status,src=PaginateOut.#.Fstatus,required
fullname=pageInfo.pageIndex,src=input.pageIndex,required
fullname=pageInfo.pageSize,src=input.pageSize,required
fullname=pageInfo.total,src=PaginateTotalOut,required
`

var schemalineIn = `
	fullname=config.openId,dst=FopenID,format=DBValidate,required
	fullname=config.type,dst=FopenIDType,enum=["1","2"],required
	fullname=config.status,dst=Fstatus,enum=["0","1"],format=number,required
 `

func TestParseInstructOut(t *testing.T) {
	instructs := jsonschemaline.ParseInstruct(schemalineOut)
	fmt.Println(instructs.String())
}
func TestParseInstructIn(t *testing.T) {
	instructs := jsonschemaline.ParseInstruct(schemalineIn)
	fmt.Println(instructs.String())
}
