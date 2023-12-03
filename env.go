package main

import (
	"os"

	"github.com/chiyoi/apricot/logs"
	"github.com/chiyoi/morph/env"
)

func init() {
	switch os.Getenv("ENV") {
	case "prod":
		env.Prod()
	default:
		env.Dev()
		logs.SetLevel(logs.LevelDebug)
	}
}
