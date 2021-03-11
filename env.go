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
	username := os.Getenv(varPrefix + "_USERNAME")
	if username == "" {
		return "", "", fmt.Errorf("%s_USERNAME is not set", varPrefix)
	}
	password := os.Getenv(varPrefix + "_PASSWORD")
	if password == "" {
		return "", "", fmt.Errorf("%s_PASSWORD is not set", varPrefix)
	}
	return username, password, nil
}

func (*Env) isThis(serverURL, serverHostname, registryDomain string) bool {
	return strings.HasSuffix(serverURL, "."+registryDomain) || (serverURL == registryDomain) ||
		strings.HasSuffix(serverHostname, "."+registryDomain) || (serverHostname == registryDomain)
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
	if e.isThis(serverURL, u.Host, "azurecr.com") {
		return e.getFor("ACR")
	}
	if e.isThis(serverURL, u.Host, "docker.io") {
		return e.getFor("DOCKER_HUB")
	}
	if e.isThis(serverURL, u.Host, "amazonaws.com") {
		return e.getFor("ECR")
	}
	if e.isThis(serverURL, u.Host, "gcr.io") {
		return e.getFor("GCR")
	}
	if e.isThis(serverURL, u.Host, "ghcr.io") {
		return e.getFor("GHCR")
	}
	if e.isThis(serverURL, u.Host, "quay.io") {
		return e.getFor("QUAY")
	}
	anyRegistryDisable, haveAnyRegistryDisable := os.LookupEnv("ANY_REGISTRY_DISABLE")
	if haveAnyRegistryDisable && anyRegistryDisable == "true" {
		return "", "", fmt.Errorf("unsupported registry: %q", serverURL)
	}
	return e.getFor("ANY_REGISTRY")
}

// List is not supported and will always error.
func (*Env) List() (map[string]string, error) {
	return nil, errors.New("list is not supported")
}
