package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/chiyoi/apricot/logs"
	"github.com/chiyoi/apricot/neko"
	"github.com/chiyoi/az"
	"github.com/chiyoi/az/cosmos"
	"github.com/chiyoi/morph/clients"
)

func main() {
	srv := &http.Server{
		Addr:    os.Getenv("ADDR"),
		Handler: Handler(),
	}

	go neko.StartServer(srv, false)
	defer neko.StopServer(srv)

	neko.Block()
}

func Handler() http.Handler {
	hostMap, err := clients.ContainerClientHostMap()
	if err != nil {
		logs.Panic(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logs.Debug("At :12381.", "r.Host:", r.Host, "r.URL:", r.URL, "r.Header:", r.Header)
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()
		var e HostMapEntry
		if err := cosmos.PointRead(ctx, hostMap, r.Host, &e); err != nil {
			if az.IsNotFound(err) {
				logs.Warning("Host not found.", "r.Host:", r.Host)
				neko.Teapot(w)
				return
			}
			logs.Error(err)
			neko.InternalServerError(w)
			return
		}

		u := r.URL
		u.Scheme = e.Target.Schema
		u.Host = e.Target.Host
		req, err := http.NewRequest(r.Method, u.String(), r.Body)
		if err != nil {
			logs.Error(err)
			neko.InternalServerError(w)
			return
		}
		for k, v := range r.Header {
			req.Header[k] = v
		}

		re, err := http.DefaultClient.Do(req)
		if err != nil {
			logs.Error(err)
			neko.InternalServerError(w)
			return
		}
		defer re.Body.Close()

		for k, v := range re.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(re.StatusCode)
		if _, err := io.Copy(w, re.Body); err != nil {
			logs.Warning(err)
		}
	})
}

type HostMapEntry struct {
	ID     string `json:"id"`
	Target struct {
		Schema string `json:"schema"`
		Host   string `json:"host"`
	} `json:"target"`
}
