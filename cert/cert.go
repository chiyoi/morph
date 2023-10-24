package cert

import (
	"github.com/chiyoi/morph/cache"
	"github.com/chiyoi/morph/clients"
	"golang.org/x/crypto/acme/autocert"
)

var (
	HostWhitelist = []string{
		"aira.neko03.moe",
		"trinity.neko03.moe",
	}
)

func Manager() (m *autocert.Manager, err error) {
	c, err := clients.BlobContainerClientCertCache()
	if err != nil {
		return
	}
	return &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      cache.NewBlobCertCache(c),
		HostPolicy: autocert.HostWhitelist(HostWhitelist...),
	}, nil
}
