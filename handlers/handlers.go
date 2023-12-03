package handlers

import (
	"net/http"

	"github.com/chiyoi/apricot/neko"
	"github.com/chiyoi/morph/handlers/registers"
)

func Root() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/ping", neko.PingHandler())
	mux.Handle("/version", neko.VersionHandler())
	mux.Handle(registers.PatternHandler("/registers/"))
	mux.Handle("/", neko.TemporaryRedirectHandler("/version"))
	return ProxyHandler(mux)
}
