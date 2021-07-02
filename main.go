package main

import (
	"context"
	"net/http"

	"github.com/luraproject/lura/logging"
	_ "github.com/ntchjb/martian-querybody/bodyquery"
	_ "github.com/ntchjb/martian-querybody/query"
)

var ClientRegisterer = registerer("martian-querybody")

type registerer string

var logger logging.Logger = nil

func (r registerer) RegisterLogger(v interface{}) {
	l, ok := v.(logging.Logger)
	if !ok {
		return
	}
	logger = l
	logger.Debug(ClientRegisterer, "client plugin loaded!!!")
}

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
}

func main() {

}
