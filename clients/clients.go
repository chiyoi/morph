package clients

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/chiyoi/az/identity"
)

var (
	EndpointAzureCosmos = os.Getenv("ENDPOINT_AZURE_COSMOS")
	Database            = os.Getenv("DATABASE")
	ContainerHostMap    = "host_map"
)

func ContainerClientHostMap() (c *azcosmos.ContainerClient, err error) {
	cred, err := identity.DefaultCredential()
	if err != nil {
		return
	}
	client, err := azcosmos.NewClient(EndpointAzureCosmos, cred, nil)
	if err != nil {
		return
	}
	return client.NewContainer(Database, ContainerHostMap)
}
