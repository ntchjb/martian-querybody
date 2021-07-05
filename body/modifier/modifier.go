package modifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ntchjb/martian-querybody/util"
)

type Config struct {
	Schema       []map[string]string               `json:"schema"`
	ValueMapping map[string]map[string]interface{} `json:"value_map"`
}

type BodyModifier struct {
	schema       []map[string]string
	valueMapping map[string]map[string]interface{}
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

func moveKey(requestBody map[string]interface{}, oldKey string, newKey string, isCopy bool) {
	mapping, key := getKeyPosition(requestBody, oldKey)

	value := mapping[key]
	if !isCopy {
		delete(mapping, key)
	}
	if newKey != "" {
		addKeyToBody(requestBody, newKey, value)
	}
}

func mapValue(requestBody map[string]interface{}, key string, valMap map[string]interface{}) {
	mapping, k := getKeyPosition(requestBody, key)
	oldValInf := mapping[k]

	oldVal := util.ConvertAnyToString(oldValInf)
	if newVal, ok := valMap[oldVal]; ok {
		mapping[k] = newVal
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
			var isCopy bool
			if len(newKey) >= 2 && newKey[0] == ',' {
				newKeySplited := strings.SplitN(newKey, ",", 2)
				newKey = newKeySplited[1]
				isCopy = true
			}

			moveKey(responseBody, oldKey, newKey, isCopy)
		}
	}

	for key, valueMap := range m.valueMapping {
		mapValue(responseBody, key, valueMap)
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
		schema:       cfg.Schema,
		valueMapping: cfg.ValueMapping,
	}, nil
}
