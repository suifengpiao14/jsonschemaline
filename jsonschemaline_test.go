package jsonschemaline_test

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/jsonschemaline"
)

var jsonStr = `{"$schema":"http://json-schema.org/draft-07/schema#","type":"object","properties":{"items":{"type":"array","allowEmpty":true,"items":{"type":"object","properties":{"id":{"type":"string","src":"hjxmba_dbCityStoreServiceConfigPaginateOut.#.Fid"},"cityId":{"type":"string","src":"hjxmba_dbCityStoreServiceConfigPaginateOut.#.FcityId"},"cityName":{"type":"string","src":"hjxmba_dbCityStoreServiceConfigPaginateOut.#.FcityName"},"status":{"type":"string","src":"hjxmba_dbCityStoreServiceConfigPaginateOut.#.Fstatus"}},"required":["id","cityId","cityName","status"]}},"pageInfo":{"type":"object","properties":{"pageIndex":{"type":"string","src":"input.pageIndex"},"pageSize":{"type":"string","format":"number","src":"input.pageSize"},"total":{"type":"string","src":"hjxmba_dbCityStoreServiceConfigTotalOut"}},"required":["pageIndex","pageSize","total"]}},"required":["items","pageInfo"]}`

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
