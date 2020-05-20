package main

import (
	"github.com/docker/docker-credential-helpers/credentials"
)

func main() {
	credentials.Serve(&Env{})
}
