package clients

import (
	"context"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/chiyoi/apricot/logs"
	"github.com/chiyoi/az/identity"
)

var (
	ContainerHostMap       = "host_map"
	ContainerSkipList      = "skip_list"
	ContainerHostWhitelist = "host_whitelist"
)

func CheckConnectivity() (ok bool) {
	db, err := databaseClient()
	if err != nil {
		logs.Warning(err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err = db.Read(ctx, nil)
	if err != nil {
		logs.Warning(err)
		return
	}

	c, err := blobClient()
	if err != nil {
		logs.Warning(err)
		return
	}
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err = c.ServiceClient().GetProperties(ctx, nil)
	if err != nil {
		logs.Warning(err)
		return
	}
	return true
}

func ContainerClientHostMap() (c *azcosmos.ContainerClient, err error) {
	client, err := databaseClient()
	if err != nil {
		return
	}
	return client.NewContainer(ContainerHostMap)
}

func ContainerClientSkipList() (c *azcosmos.ContainerClient, err error) {
	client, err := databaseClient()
	if err != nil {
		return
	}
	return client.NewContainer(ContainerSkipList)
}

func ContainerClientHostWhitelist() (c *azcosmos.ContainerClient, err error) {
	client, err := databaseClient()
	if err != nil {
		return
	}
	return client.NewContainer(ContainerHostWhitelist)
}

func BlobContainerClientCertCache() (c *container.Client, err error) {
	var (
		BlobContainerCertCache = os.Getenv("BLOB_CONTAINER_CERT_CACHE")
	)
	client, err := blobClient()
	return client.ServiceClient().NewContainerClient(BlobContainerCertCache), err
}

func databaseClient() (c *azcosmos.DatabaseClient, err error) {
	var (
		EndpointAzureCosmos = os.Getenv("ENDPOINT_AZURE_COSMOS")
		Database            = os.Getenv("DATABASE")
	)

	cred, err := identity.DefaultCredential()
	if err != nil {
		return
	}

	client, err := azcosmos.NewClient(EndpointAzureCosmos, cred, nil)
	if err != nil {
		return
	}
	return client.NewDatabase(Database)
}

func blobClient() (c *azblob.Client, err error) {
	var (
		EndpointAzureBlob = os.Getenv("ENDPOINT_AZURE_BLOB")
	)

	cred, err := identity.DefaultCredential()
	if err != nil {
		return
	}
	return azblob.NewClient(EndpointAzureBlob, cred, nil)
}
