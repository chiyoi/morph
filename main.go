package main

import (
	"net/http"

	"github.com/chiyoi/apricot/neko"
	"github.com/chiyoi/morph/env"
	"github.com/chiyoi/morph/handlers"
)

func main() {
	srv := &http.Server{
		Addr:    env.Addr,
		Handler: handlers.Root(),
	}

	go neko.StartServer(srv, false)
	defer neko.StopServer(srv)

	neko.Block()
}
