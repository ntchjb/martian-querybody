package modifier

import (
	"encoding/json"
	"net/http"

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

// ModifyRequest converts body to query parameters based on given key name mapping
func (m *BodyQueryModifier) ModifyRequest(req *http.Request) error {
	query := req.URL.Query()
	for oldKey, newKey := range m.keyMapping {
		if oldKey == newKey {
			continue
		}
		if values := query[oldKey]; len(values) > 0 {
			query.Del(oldKey)
			for _, value := range values {
				query.Add(newKey, value)
			}
		}
	}
	for newKey, newValue := range m.valueMapping {
		if values := query[newKey]; len(values) > 0 {
			query[newKey][newValue.Index] = util.ConvertAnyToString(newValue.Value)
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
