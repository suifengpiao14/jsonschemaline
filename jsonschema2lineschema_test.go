package jsonschemaline_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/suifengpiao14/jsonschemaline"
)

func TestJsonSchema2LineSchema(t *testing.T) {
	var jsonschmaStr = `{"$schema":"http://json-schema.org/draft-07/schema#","$id":"execAPIInquiryScreenIdentifyUpdate","type":"object","required":["config"],"properties":{"config":{"properties":{"id":{"type":"string","format":"number"},"status":{"type":"string","enum":[1,"2"]},"identify":{"type":"string","minLength":1},"merchantId":{"type":"string","format":"number"},"merchantName":{"type":"string","minLength":1},"operateName":{"type":"string","minLength":1},"storeId":{"type":"string","format":"number"},"storeName":{"type":"string","minLength":1}},"type":"object","required":["id","status","identify","merchantId","merchantName","operateName","storeId","storeName"]}}}`

	lineschemaStr, err := jsonschemaline.JsonSchema2LineSchema(jsonschmaStr)
	if err != nil {
		panic(err)
	}

	_ = lineschemaStr
	fmt.Println(lineschemaStr)

}

func parseJSONSchema(schema map[string]interface{}, prefix string) []string {
	output := []string{}
	for key, value := range schema {
		fieldName := fmt.Sprintf("%s%s", prefix, key)
		switch valueType := value.(type) {
		case map[string]interface{}:
			// 递归处理子对象
			output = append(output, parseJSONSchema(valueType, fieldName+".")...)
		default:
			// 将键值对格式化为字符串
			output = append(output, fmt.Sprintf("%s=%v", fieldName, value))
		}
	}
	return output
}

func TestParseInput(t *testing.T) {
	// 示例 JSON Schema
	jsonSchema := `{
		"type": "object",
		"properties": {
			"config": {
				"type": "array",
				"items": {
					"type":"object",
					"properties":{
						"id": {
							"type": "string",
							"format": "number"
						},
						"status": {
							"type": "string",
							"enum": ["1", "2"]
						},
						"identify": {
							"type": "string",
							"minLength": 1
						},
						"merchantId": {
							"type": "string",
							"format": "number"
						},
						"merchantName": {
							"type": "string",
							"minLength": 1
						},
						"operateName": {
							"type": "string",
							"minLength": 1
						},
						"storeId": {
							"type": "string",
							"format": "number"
						},
						"storeName": {
							"type": "string",
							"minLength": 1
						}
					}
				},
				"required": ["id", "status", "identify", "merchantId", "merchantName", "operateName", "storeId", "storeName"]
			}
		}
	}`

	var schema map[string]interface{}
	err := json.Unmarshal([]byte(jsonSchema), &schema)
	if err != nil {
		panic(err)
	}

	output := parseJSONSchema(schema, "")

	// 输出键值对格式
	fmt.Println(strings.Join(output, "\n"))
}
