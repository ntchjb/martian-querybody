package modifier

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Config struct {
	ValueMapping map[string]NewValue `json:"value_map"`
	KeyMapping   map[string]string   `json:"key_map"`
}

type NewValue struct {
	Index int         `json:"index"`
	Value interface{} `json:"value"`
}

type QueryBodyModifier struct {
	keyMapping   map[string]string
	valueMapping map[string]NewValue
}

func convertAnyToString(v interface{}) string {
	switch result := v.(type) {
	case bool:
		return strconv.FormatBool(result)
	case float64:
		return strconv.FormatFloat(result, 'f', -1, 64)
	case string:
		return result
	default:
		return ""
	}
}

func getValueFromKeyInBody(requestBody map[string]interface{}, key string) []string {
	keys := strings.Split(key, ".")
	currentValue := requestBody
	var finalResult []string

loop:
	for _, k := range keys {
		value := currentValue[k]

		switch result := value.(type) {
		case map[string]interface{}:
			currentValue = result
		case []interface{}:
			for _, el := range result {
				finalResult = append(finalResult, convertAnyToString(el))
			}
			break loop
		default:
			finalResult = append(finalResult, convertAnyToString(result))
			break loop
		}
	}

	return finalResult
}

// ModifyRequest converts body to query parameters based on given key name mapping
func (m *QueryBodyModifier) ModifyRequest(req *http.Request) error {
	decoder := json.NewDecoder(req.Body)
	var requestBody map[string]interface{}
	err := decoder.Decode(&requestBody)
	if err != nil {
		return fmt.Errorf("unable to parse request body: %w", err)
	}
	query := req.URL.Query()
	for oldKey, newKey := range m.keyMapping {
		values := getValueFromKeyInBody(requestBody, oldKey)
		for i, value := range values {
			if value != "" {
				if newVal, ok := m.valueMapping[newKey]; ok && newVal.Index == i {
					query.Add(newKey, convertAnyToString(newVal.Value))
					continue
				}
				query.Add(newKey, value)
			}
		}
	}
	req.URL.RawQuery = query.Encode()

	return nil
}

func FromJSON(b []byte) (*QueryBodyModifier, error) {
	cfg := &Config{}
	if err := json.Unmarshal(b, cfg); err != nil {
		return nil, err
	}

	return &QueryBodyModifier{
		keyMapping:   cfg.KeyMapping,
		valueMapping: cfg.ValueMapping,
	}, nil
}
