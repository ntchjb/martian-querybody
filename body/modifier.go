package body

import (
	"github.com/google/martian/parse"
	"github.com/ntchjb/martian-querybody/body/modifier"
)

func init() {
	parse.Register("body.fromBody", FromJSON)
}

func FromJSON(b []byte) (*parse.Result, error) {
	msg, err := modifier.FromJSON(b)
	if err != nil {
		return nil, err
	}

	return parse.NewResult(msg, []parse.ModifierType{parse.Response})
}
