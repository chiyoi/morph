package main

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/chiyoi/apricot/logs"
	"github.com/chiyoi/apricot/neko"
	"github.com/chiyoi/az"
	"github.com/chiyoi/az/cosmos"
	"github.com/chiyoi/morph/cert"
	"github.com/chiyoi/morph/clients"
)

var (
	SkipHost = map[string]bool{
		"morph.neko03.moe": true,
		"localhost:12380":  true,
	}
)

func main() {
	m, err := cert.Manager()
	if err != nil {
		logs.Panic(err)
	}

	if os.Getenv("ENV") == "prod" {
		HTTPSSrv := &http.Server{
			Addr:      ":https",
			Handler:   RootHandler(),
			TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		}
		HTTPSrv := &http.Server{
			Addr:    ":http",
			Handler: m.HTTPHandler(RedirectToHTTPSHandler()),
		}

		go neko.StartServer(HTTPSrv, false)
		defer neko.StopServer(HTTPSrv)
		go neko.StartServer(HTTPSSrv, true)
		defer neko.StopServer(HTTPSSrv)
	} else {
		srv := &http.Server{
			Addr:    os.Getenv("ADDR"),
			Handler: RootHandler(),
		}

		go neko.StartServer(srv, false)
		defer neko.StopServer(srv)
	}

	neko.Block()
}

func RootHandler() http.Handler {
	hostMap, err := clients.ContainerClientHostMap()
	if err != nil {
		logs.Panic(err)
	}
	skipList, err := clients.ContainerClientSkipList()
	if err != nil {
		logs.Panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/ping", neko.PingHandler())
	mux.Handle("/version", neko.VersionHandler())
	mux.Handle("/readiness", neko.ReadinessHandler(clients.CheckConnectivity))
	mux.Handle("/", neko.TeapotHandler())
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logs.Info("Inbound received.", "r.Host:", r.Host, "r.URL:", r.URL)
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()
		exist, err := cosmos.PointKeyExist(ctx, skipList, r.Host)
		if err != nil {
			logs.Error(err)
			neko.InternalServerError(w)
			return
		}
		if exist {
			logs.Info("Skipped.")
			mux.ServeHTTP(w, r)
			return
		}

		logs.Info("Lookup host.")
		ctx, cancel = context.WithTimeout(r.Context(), time.Second*10)
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
			logs.Warning(err)
		}
	})
}

func RedirectToHTTPSHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := r.URL
		u.Scheme = "https"
		u.Host = r.Host
		neko.TemporaryRedirect(w, r, u.String())
	})
}

type HostMapEntry struct {
	ID     string `json:"id"`
	Target struct {
		Schema string `json:"schema"`
		Host   string `json:"host"`
	} `json:"target"`
}
