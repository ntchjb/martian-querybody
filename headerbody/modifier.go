package headerbody

import (
	"github.com/google/martian/parse"
	"github.com/ntchjb/martian-querybody/headerbody/modifier"
)

func init() {
	parse.Register("body.fromHeader", FromJSON)
}

func FromJSON(b []byte) (*parse.Result, error) {
	msg, err := modifier.FromJSON(b)
	if err != nil {
		return nil, err
	}

	return parse.NewResult(msg, []parse.ModifierType{parse.Response})
}
