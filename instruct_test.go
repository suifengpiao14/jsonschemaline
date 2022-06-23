package jsonschemaline_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/jsonschemaline"
)

var schemalineOut = `
version=http://json-schema.org/draft-07/schema#,id=response
fullname=items[].id,src=PaginateOut.#.Fid,required
fullname=items[].openId,src=PaginateOut.#.Fopen_id,required
fullname=items[].type,src=PaginateOut.#.Fopen_id_type,required
fullname=items[].status,src=PaginateOut.#.Fstatus,required
fullname=pageInfo.pageIndex,src=input.pageIndex,required
fullname=pageInfo.pageSize,src=input.pageSize,required
fullname=pageInfo.total,src=PaginateTotalOut,required
`

var schemalineOut2 = `
version=http://json-schema.org/draft-07/schema#,id=mainOut
fullname=items[].id,required,src=PaginateOut.#.Fid
fullname=items[].identify,required,src=PaginateOut.#.Fidentify
fullname=items[].merchantId,required,src=PaginateOut.#.Fmerchant_id
fullname=items[].merchantName,required,src=PaginateOut.#.Fmerchant_name
fullname=items[].operateName,required,src=PaginateOut.#.Foperate_name
fullname=items[].status,required,src=PaginateOut.#.Fstatus
fullname=items[].storeId,required,src=PaginateOut.#.Fstore_id
fullname=items[].storeName,required,src=PaginateOut.#.Fstore_name
fullname=items[].createTime,required,src=PaginateOut.#.Fcreate_time
fullname=items[].updateTime,required,src=PaginateOut.#.Fupdate_time
fullname=pageInfo.pageIndex,required,src=mainIn.pageIndex
fullname=pageInfo.pageSize,required,src=mainIn.pageSize
fullname=pageInfo.total,required,src=PaginateTotalOut
`

var schemalineIn = `
    version=http://json-schema.org/draft-07/schema#,id=request
	fullname=config.openId,dst=FopenID,format=DBValidate,required
	fullname=config.type,dst=FopenIDType,enum=["1","2"],required
	fullname=config.status,dst=Fstatus,enum=["0","1"],format=number,required
 `
var schemalineIn2 = `
    version=http://json-schema.org/draft-07/schema#,id=mainIn
	fullname=pageSize,format=number,required,dst=Limit
fullname=pageIndex,format=number,required,dst={{setValue . "Offset" (mul  (getValue .  "mainIn.pageIndex")   (getValue . "mainIn.pageSize"))}}
 `

func TestParseInstructOut(t *testing.T) {
	//instructs := jsonschemaline.ParseOneInstructTp(schemalineOut)
	instructTps := jsonschemaline.ParseInstructTp(schemalineOut2)
	for _, instructTp := range instructTps {
		fmt.Println(instructTp.String())
	}
}
func TestParseInstructIn(t *testing.T) {
	instructTps := jsonschemaline.ParseInstructTp(schemalineIn2)
	for _, instructTp := range instructTps {
		fmt.Println(instructTp.String())
	}
}
