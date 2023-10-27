package cert

import (
	"context"
	"fmt"
	"time"

	"github.com/chiyoi/az/cosmos"
	"github.com/chiyoi/morph/cache"
	"github.com/chiyoi/morph/clients"
	"golang.org/x/crypto/acme/autocert"
)

func Manager() (m *autocert.Manager, err error) {
	hwl, err := clients.ContainerClientHostWhitelist()
	if err != nil {
		return
	}
	c, err := clients.BlobContainerClientCertCache()
	if err != nil {
		return
	}
	return &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  cache.NewBlobCertCache(c),
		HostPolicy: func(ctx context.Context, host string) (err error) {
			ctx, cancel := context.WithTimeout(ctx, time.Second*10)
			defer cancel()
			exist, err := cosmos.PointKeyExist(ctx, hwl, host)
			if err != nil {
				return
			}
			if !exist {
				return fmt.Errorf("unrecognised host (%s)", host)
			}
			return
		},
	}, nil
}
