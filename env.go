package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/docker/docker-credential-helpers/credentials"
)

// Env handles secrets using environment variables (read-only)
type Env struct{}

// Add is not supported and will always error.
func (*Env) Add(*credentials.Credentials) error {
	return errors.New("add is not supported")
}

// Add is not supported and will always error.
func (*Env) Delete(string) error {
	return errors.New("delete is not supported")
}

func (*Env) getFor(varPrefix string) (string, string, error) {
	publicAccess, havePublicAccess := os.LookupEnv(varPrefix + "_PUBLIC_ACCESS_ONLY")
	if havePublicAccess && publicAccess == "true" {
		return "", "", nil
	}
	username, haveUsername := os.LookupEnv(varPrefix + "_USERNAME")
	if !haveUsername {
		return "", "", fmt.Errorf("%s_USERNAME is not set", varPrefix)
	}
	password, havePassword := os.LookupEnv(varPrefix + "_PASSWORD")
	if !havePassword {
		return "", "", fmt.Errorf("%s_USERNAME is not set", varPrefix)
	}
	return username, password, nil
}

// Get returns the username and secret to use for a given registry server URL.
func (e *Env) Get(serverURL string) (string, string, error) {
	if serverURL == "" {
		return "", "", errors.New("missing server URL")
	}

	u, err := url.Parse(serverURL)
	if err != nil {
		return "", "", err
	}
	if strings.HasSuffix(u.Host, ".docker.io") {
		return e.getFor("DOCKER_HUB")
	}
	if strings.HasSuffix(u.Host, ".quay.io") {
		return e.getFor("QUAY")
	}
	return e.getFor("ANY_REGISTRY")
}

// List is not supported and will always error.
func (*Env) List() (map[string]string, error) {
	return nil, errors.New("list is not supported")
}
