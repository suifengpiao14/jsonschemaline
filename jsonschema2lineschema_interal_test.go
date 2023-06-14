package jsonschemaline

import (
	"fmt"
	"testing"

	"github.com/suifengpiao14/kvstruct"
)

func TestDealRequired(t *testing.T) {

	kvs := kvstruct.KVS{
		{Key: "config.required.1", Value: "status"},
		{Key: "config.type", Value: "object"},
		{Key: "config.required", Value: "true"},
	}
	newKvs := dealRequired(kvs)
	fmt.Println(newKvs)

}
