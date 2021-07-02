package bodyquery

import (
	"github.com/google/martian/parse"
	"github.com/ntchjb/martian-querybody/query/modifier"
)

func init() {
	parse.Register("query.fromQuery", FromJSON)
}

func FromJSON(b []byte) (*parse.Result, error) {
	msg, err := modifier.FromJSON(b)
	if err != nil {
		return nil, err
	}

	return parse.NewResult(msg, []parse.ModifierType{parse.Request})
}
