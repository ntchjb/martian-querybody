package modifier

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertAnyToString(t *testing.T) {
	var v interface{}

	v = float64(15)
	require.Equal(t, "15", convertAnyToString(v))

	v = "hello"
	require.Equal(t, "hello", convertAnyToString(v))

	v = true
	require.Equal(t, "true", convertAnyToString(v))

	v = float64(12.345)
	require.Equal(t, "12.345", convertAnyToString(v))
}

func TestGetValueFromKeyInBody(t *testing.T) {
	requestBody := make(map[string]interface{})
	ccc := make(map[string]interface{})
	ccc["cca"] = "testdeep"
	ccc["ccb"] = []interface{}{"testdeep1", "testdeep2"}
	cc := make(map[string]interface{})
	cc["ca"] = "testca"
	cc["cb"] = ccc
	requestBody["a"] = true
	requestBody["b"] = []interface{}{"ba", 1.25, true}
	requestBody["c"] = cc

	require.Equal(t, []string{"true"}, getValueFromKeyInBody(requestBody, "a"))
	require.Equal(t, []string{"ba", "1.25", "true"}, getValueFromKeyInBody(requestBody, "b"))
	require.Equal(t, []string{"testca"}, getValueFromKeyInBody(requestBody, "c.ca"))
	require.Equal(t, []string{"testdeep"}, getValueFromKeyInBody(requestBody, "c.cb.cca"))
	require.Equal(t, []string{"true"}, getValueFromKeyInBody(requestBody, "a"))
	require.Equal(t, []string{"testdeep1", "testdeep2"}, getValueFromKeyInBody(requestBody, "c.cb.ccb"))
}

func TestModifyRequest(t *testing.T) {
	request, err := http.NewRequest(http.MethodPost, "/test/url", bytes.NewBuffer([]byte(`
	{
		"key1": {
			"key2": "valueKey2"
		},
		"key2": [
			13,
			14.5,
			16
		],
		"key3": "valueKey3"
	}
	`)))
	require.NoError(t, err)

	cfg := `
	{
		"key_mapping": {
			"key1.key2": "key12",
			"key2": "key2",
			"key3": "newKey3"
		}
	}
	`
	modifier, err := FromJSON([]byte(cfg))
	require.NoError(t, err)

	err = modifier.ModifyRequest(request)
	require.NoError(t, err)
	query := request.URL.Query()
	require.Equal(t, "valueKey2", query["key12"][0])
	require.Equal(t, []string{"13", "14.5", "16"}, query["key2"])
	require.Equal(t, "valueKey3", query["newKey3"][0])
}
