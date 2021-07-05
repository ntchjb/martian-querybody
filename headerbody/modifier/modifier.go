package modifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Config struct {
	KeyMapping map[string]string `json:"key_map"`
}

type HeaderBodyModifier struct {
	keyMapping map[string]string
}

// ModifyRequest converts body to query parameters based on given key name mapping
func (m *HeaderBodyModifier) ModifyResponse(res *http.Response) error {
	decoder := json.NewDecoder(res.Body)
	var requestBody map[string]interface{}
	err := decoder.Decode(&requestBody)
	if err != nil {
		return fmt.Errorf("unable to parse request body: %w", err)
	}

	for headerKey, bodyKey := range m.keyMapping {
		keys := strings.Split(bodyKey, ".")
		obj := requestBody
		var selectedKey string

	loop:
		for i, key := range keys {
			if i == len(keys)-1 {
				selectedKey = key
				break
			}
			switch result := obj[key].(type) {
			case map[string]interface{}:
				obj = result
			default:
				selectedKey = key
				break loop
			}
		}
		value := res.Header.Get(headerKey)
		if value != "" {
			obj[selectedKey] = value
		}
	}

	newResponse, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("unable to parse response after modified: %w", err)
	}

	res.Body.Close()
	res.Body = ioutil.NopCloser(bytes.NewBuffer(newResponse))
	return nil
}

func FromJSON(b []byte) (*HeaderBodyModifier, error) {
	cfg := &Config{}
	if err := json.Unmarshal(b, cfg); err != nil {
		return nil, err
	}

	return &HeaderBodyModifier{
		keyMapping: cfg.KeyMapping,
	}, nil
}
