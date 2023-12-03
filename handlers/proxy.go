package handlers

import (
	"io"
	"net/http"

	"github.com/chiyoi/apricot/logs"
	"github.com/chiyoi/apricot/neko"
	"github.com/chiyoi/az"
	"github.com/chiyoi/az/cosmos"
	"github.com/chiyoi/iter/res"
	"github.com/chiyoi/morph/containers"
	"github.com/chiyoi/morph/env"
)

func ProxyHandler(fallback http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logs.Info("Inbound received.", "r.Host:", r.Host, "r.URL:", r.URL)
		if env.FallbackToSelf(r.Host) {
			logs.Info("Fallback to self.")
			fallback.ServeHTTP(w, r)
			return
		}

		logs.Info("Lookup host.")
		var target struct {
			Schema string `json:"schema"`
			Host   string `json:"host"`
		}
		hm, err := containers.Client(containers.HostMap)
		err = res.C(hm, err, cosmos.PointRead(ctx, r.Host, &target))
		switch {
		case az.IsNotFound(err):
			logs.Warning("Host not found.", "r.Host:", r.Host)
			neko.Teapot(w)
			return
		case err != nil:
			logs.Error(err)
			neko.InternalServerError(w)
			return
		}

		u := r.URL
		u.Scheme = target.Schema
		u.Host = target.Host
		req, err := http.NewRequest(r.Method, u.String(), r.Body)
		if err != nil {
			logs.Error(err)
			neko.InternalServerError(w)
			return
		}
		for k, v := range r.Header {
			req.Header[k] = v
		}

		logs.Info("Request target.", "u:", u)
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
			logs.Error(err)
		}
	})
}
