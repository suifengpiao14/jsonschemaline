package jsonschemaline

import (
	"fmt"
	"testing"
)

func TestGetTokens(t *testing.T) {
	tokens := getTokens()
	fmt.Println(tokens)
}

func TestParserLine(t *testing.T) {
	line := `
	fullname=pageSize,format=number,enum=[1,3,3],required,dst=Limit,example={"id":1,"version":2},comments=枚举值(1-OK,0-false)
	`
	ParserLine(line)
}
