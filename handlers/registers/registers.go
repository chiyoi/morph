package registers

import (
	"net/http"
	"strings"

	"github.com/chiyoi/apricot/neko"
	"github.com/chiyoi/morph/handlers/registers/host"
)

func PatternHandler(pattern string) (string, http.Handler) {
	return pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := neko.TrimPattern(r.URL.Path, pattern)
		if len(p) != 0 {
			ps := strings.Split(p, "/")
			hostName := ps[0]
			host.DynamicHandler(hostName, ps[1:]).ServeHTTP(w, r)
			return
		}
	})
}
