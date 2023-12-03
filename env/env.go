package env

import "os"

const (
	EndpointCosmos = "https://neko03cosmos.documents.azure.com:443/"
)

var (
	Addr           string
	Database       string
	FallbackToSelf func(string) bool
)

func Dev() {
	os.Setenv("VERSION", "Morph v0.3.0-dev")
	Addr = ":12380"
	Database = "neko0001"
	FallbackToSelf = func(s string) bool {
		switch s {
		case "booklet.local:12380", "localhost:12380":
			return true
		}
		return false
	}
}

func Prod() {
	os.Setenv("VERSION", "Morph v0.3.0")
	Addr = ":http"
	Database = "trinity"
	FallbackToSelf = func(s string) bool {
		switch s {
		case "morph.neko03.moe", "localhost":
			return true
		}
		return false
	}
}
