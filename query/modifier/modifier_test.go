package modifier

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModifyRequest(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "/test/url", nil)
	require.NoError(t, err)

	query := make(url.Values)
	query.Add("key1", "valueKey1")
	query.Add("key2", "valueKey2")
	query.Add("key2", "valueKey22")
	query.Add("key3", "valueKey3")
	query.Add("key4", "valueKey4")
	query.Add("key4", "valueKey44")
	request.URL.RawQuery = query.Encode()

	cfg := `
	{
		"key_map": {
			"key1": "newKey1",
			"key2": "newKey2",
			"key4": "newKey4"
		},
		"value_map": {
			"key3": {
				"value": "newValueKey3"
			},
			"newKey4": {
				"value": "newValueKey4",
				"index": 1
			}
		}
	}
	`
	modifier, err := FromJSON([]byte(cfg))
	require.NoError(t, err)
	err = modifier.ModifyRequest(request)
	require.NoError(t, err)

	query = request.URL.Query()
	require.Equal(t, "valueKey1", query["newKey1"][0])
	require.Equal(t, []string{"valueKey2", "valueKey22"}, query["newKey2"])
	require.Equal(t, "newValueKey3", query["key3"][0])
	require.Equal(t, []string{"valueKey4", "newValueKey4"}, query["newKey4"])
}
