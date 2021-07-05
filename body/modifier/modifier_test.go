package modifier

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModifyResponse(t *testing.T) {
	response := http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body: ioutil.NopCloser(bytes.NewBufferString(`
		{
			"key1": "valueKey1",
			"key2": {
				"key22": "valueKey22"
			},
			"key3": "valueKey3",
			"key4": "valueKey4",
			"key5": {
				"key55": "valueKey55"
			},
			"key6": "valueKey6",
			"key7": "valueKey7"
		}`)),
	}

	defer response.Body.Close()

	cfg := `
	{
		"schema": [
			{
				"key1": "key1.key2",
				"key3": "key33.key3",
				"key4": "key33.key4",
				"key5.key55": "key55",
				"key6": "key66",
				"key7": ",key77"
			}, {
				"key5": ""
			}
		]
	}`
	modifier, err := FromJSON([]byte(cfg))
	require.NoError(t, err)

	err = modifier.ModifyResponse(&response)
	require.NoError(t, err)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, "{\"key1\":{\"key2\":\"valueKey1\"},\"key2\":{\"key22\":\"valueKey22\"},\"key33\":{\"key3\":\"valueKey3\",\"key4\":\"valueKey4\"},\"key55\":\"valueKey55\",\"key66\":\"valueKey6\",\"key7\":\"valueKey7\",\"key77\":\"valueKey7\"}", string(bodyBytes))
}
