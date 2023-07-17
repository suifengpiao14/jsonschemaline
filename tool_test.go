package jsonschemaline_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/jsonschemaline"
)

func TestJsonMerge(t *testing.T) {
	t.Run("case1", func(t *testing.T) {
		defaultJson := `{"pageSize":"20","remark":"hello world"}`
		specialJson := `{"pageIndex":"0","pageSize":""}`

		merge, err := jsonschemaline.MergeDefault(specialJson, defaultJson)
		require.NoError(t, err)
		expected := `{"pageIndex":"0","pageSize":"20","remark":"hello world"}`
		assert.Equal(t, expected, merge)
	})

	t.Run("case2", func(t *testing.T) {
		defaultJson := "{\"channelId\":\"10000012\",\"selectors\":\"\",\"sn\":\"\"}"
		specialJson := "{\"channelId\":\"0\",\"selectors\":[{\"id\":\"\",\"type\":\"imei\"}],\"sn\":\"2023071710052946420\",\"productId\":\"66140\"}"
		merge, err := jsonschemaline.MergeDefault(specialJson, defaultJson)
		require.NoError(t, err)
		expected := `{"channelId":"0","selectors":[{"id":"","type":"imei"}],"sn":"2023071710052946420","productId":"66140"}`
		assert.Equal(t, expected, merge)
	})

}
