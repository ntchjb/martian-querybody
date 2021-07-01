package main

import (
	"github.com/devopsfaith/krakend-martian/register"
	"github.com/ntchjb/martian-querybody/querybody/modifier"
)

func init() {
	register.Set("query.fromJSONBody", []register.Scope{register.ScopeRequest}, func(b []byte) (interface{}, error) {
		return modifier.FromJSON(b)
	})
}

func main() {

}
