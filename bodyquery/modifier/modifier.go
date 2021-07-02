package modifier

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ntchjb/martian-querybody/util"
)

type Config struct {
	ValueMapping map[string]NewValue `json:"value_map"`
	KeyMapping   map[string]string   `json:"key_map"`
}

type NewValue struct {
	Index int         `json:"index"`
	Value interface{} `json:"value"`
}

type BodyQueryModifier struct {
	keyMapping   map[string]string
	valueMapping map[string]NewValue
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
				finalResult = append(finalResult, util.ConvertAnyToString(el))
			}
			break loop
		default:
			finalResult = append(finalResult, util.ConvertAnyToString(result))
			break loop
		}
	}

	return finalResult
}

// ModifyRequest converts body to query parameters based on given key name mapping
func (m *BodyQueryModifier) ModifyRequest(req *http.Request) error {
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
					query.Add(newKey, util.ConvertAnyToString(newVal.Value))
					continue
				}
				query.Add(newKey, value)
			}
		}
	}
	req.URL.RawQuery = query.Encode()

	return nil
}

func FromJSON(b []byte) (*BodyQueryModifier, error) {
	cfg := &Config{}
	if err := json.Unmarshal(b, cfg); err != nil {
		return nil, err
	}

	return &BodyQueryModifier{
		keyMapping:   cfg.KeyMapping,
		valueMapping: cfg.ValueMapping,
	}, nil
}
