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
	Schema []map[string]string `json:"schema"`
}

type BodyModifier struct {
	schema []map[string]string
}

func getKeyPosition(m map[string]interface{}, key string) (map[string]interface{}, string) {
	keys := strings.Split(key, ".")
	currentValue := m
	selectedValue := m
	selectedKey := ""

loop:
	for _, k := range keys {
		selectedValue = currentValue
		selectedKey = k
		switch result := currentValue[k].(type) {
		case map[string]interface{}:
			currentValue = result
		default:
			break loop
		}
	}

	return selectedValue, selectedKey
}

func moveKey(requestBody map[string]interface{}, oldKey string, newKey string) {
	mapping, key := getKeyPosition(requestBody, oldKey)

	value := mapping[key]
	delete(mapping, key)
	if newKey != "" {
		addKeyToBody(requestBody, newKey, value)
	}
}

func addKeyToBody(m map[string]interface{}, key string, value interface{}) {
	keys := strings.Split(key, ".")
	currentValue := m
	selectedValue := m
	selectedKey := ""

loop:
	for i, k := range keys {
		selectedValue = currentValue
		selectedKey = k
		switch result := currentValue[k].(type) {
		case map[string]interface{}:
			currentValue = result
		case nil:
			if i >= len(keys)-1 {
				break loop
			}
			newMap := make(map[string]interface{})
			currentValue[k] = newMap
			currentValue = newMap
		default:
			break loop
		}
	}

	selectedValue[selectedKey] = value
}

// ModifyResponse converts body new body based on given key name mapping
func (m *BodyModifier) ModifyResponse(res *http.Response) error {
	decoder := json.NewDecoder(res.Body)
	var responseBody map[string]interface{}
	err := decoder.Decode(&responseBody)
	if err != nil {
		return fmt.Errorf("unable to parse request body: %w", err)
	}

	for _, orderedGroup := range m.schema {
		for oldKey, newKey := range orderedGroup {
			moveKey(responseBody, oldKey, newKey)
		}
	}

	newResponse, err := json.Marshal(responseBody)
	if err != nil {
		return fmt.Errorf("unable to parse response after modified: %w", err)
	}

	res.Body.Close()
	res.Body = ioutil.NopCloser(bytes.NewBuffer(newResponse))
	return nil
}

func FromJSON(b []byte) (*BodyModifier, error) {
	cfg := &Config{}
	if err := json.Unmarshal(b, cfg); err != nil {
		return nil, err
	}

	return &BodyModifier{
		schema: cfg.Schema,
	}, nil
}
