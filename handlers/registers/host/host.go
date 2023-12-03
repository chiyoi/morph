package host

import (
	"net/http"

	"github.com/chiyoi/apricot/kitsune"
	"github.com/chiyoi/apricot/logs"
	"github.com/chiyoi/apricot/neko"
	"github.com/chiyoi/az"
	"github.com/chiyoi/az/cosmos"
	"github.com/chiyoi/iter/res"
	"github.com/chiyoi/morph/containers"
)

func DynamicHandler(host string, ps []string) http.Handler {
	handleGet := func(w http.ResponseWriter, r *http.Request) {
		logs.Info("Get register.")
		ctx := r.Context()
		var target Target
		c, err := containers.Client(containers.HostMap)
		err = res.C(c, err, cosmos.PointRead(ctx, host, &target))
		if az.IsNotFound(err) {
			logs.Warning("Host not found:", host)
			neko.Teapot(w)
			return
		}
		if err != nil {
			logs.Error(err)
			neko.InternalServerError(w)
		}
		kitsune.Respond(w, target)
	}

	handlePut := func(w http.ResponseWriter, r *http.Request) {
		logs.Info("Put register.")
		ctx := r.Context()
		var target Target
		kitsune.ParseRequest(r, &target)
		c, err := containers.Client(containers.HostMap)
		err = res.C(c, err, cosmos.PointUpsert(ctx, host, target))
		if err != nil {
			logs.Error(err)
			neko.InternalServerError(w)
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(ps) != 0 {
			logs.Warning("Redundant paths:", ps)
			neko.BadRequest(w)
			return
		}
		switch r.Method {
		case http.MethodGet:
			handleGet(w, r)
		case http.MethodPut:
			handlePut(w, r)
		default:
			logs.Warning("Method not allowed:", r.Method)
			neko.MethodNotAllowed(w)
			return
		}
	})
}

type Target struct {
	Schema string `json:"schema"`
	Host   string `json:"host"`
}
