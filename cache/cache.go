package cache

import (
	"context"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/chiyoi/az"
	"golang.org/x/crypto/acme/autocert"
)

type BlobCertCache struct {
	c *container.Client
}

var _ autocert.Cache = (*BlobCertCache)(nil)

func NewBlobCertCache(c *container.Client) *BlobCertCache {
	return &BlobCertCache{c}
}

// Delete implements autocert.Cache.
func (bcc *BlobCertCache) Delete(ctx context.Context, key string) (err error) {
	if _, err := bcc.c.NewBlockBlobClient(key).Delete(ctx, nil); err != nil {
		if az.IsNotFound(err) {
			return nil
		}
		return err
	}
	return
}

// Get implements autocert.Cache.
func (bcc *BlobCertCache) Get(ctx context.Context, key string) (bs []byte, err error) {
	re, err := bcc.c.NewBlockBlobClient(key).DownloadStream(ctx, nil)
	if err != nil {
		if az.IsNotFound(err) {
			return nil, autocert.ErrCacheMiss
		}
		return
	}
	return io.ReadAll(re.Body)
}

// Put implements autocert.Cache.
func (bcc *BlobCertCache) Put(ctx context.Context, key string, data []byte) (err error) {
	_, err = bcc.c.NewBlockBlobClient(key).UploadBuffer(ctx, data, nil)
	return
}
