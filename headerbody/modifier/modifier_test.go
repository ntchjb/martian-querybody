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
			"key1": {
				"key2": "valueKey2"
			},
			"key2": [
				13,
				16
			],
			"key3": "valueKey3"
		}`)),
	}
	response.Header.Add("header1", "valueHeader1")
	response.Header.Add("header2", "valueHeader21")
	response.Header.Add("header2", "valueHeader22")

	defer response.Body.Close()

	cfg := `
	{
		"key_map": {
			"header1": "key1.key2",
			"header2": "key6",
			"header3": "key3"
		}
	}`
	modifier, err := FromJSON([]byte(cfg))
	require.NoError(t, err)

	err = modifier.ModifyResponse(&response)
	require.NoError(t, err)

	bodyBytes, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, string(bodyBytes), "{\"key1\":{\"key2\":\"valueHeader1\"},\"key2\":[13,16],\"key3\":\"valueKey3\",\"key6\":\"valueHeader21\"}")
}
