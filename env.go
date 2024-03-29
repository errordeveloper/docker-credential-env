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

func (e *Env) getForKnownRegistry(serverURL string, u *url.URL) (string, string, error) {
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
	return "", "", fmt.Errorf("unsupported registry %s", serverURL)
}

// Get returns the username and secret to use for a given registry server URL.
func (e *Env) Get(serverURL string) (string, string, error) {
	if serverURL == "" {
		return "", "", errors.New("missing server URL")
	}

	parsedServerURL, err := url.Parse(serverURL)
	if err != nil {
		return "", "", err
	}

	username, password, err := e.getForKnownRegistry(serverURL, parsedServerURL)
	if err == nil {
		return username, password, nil
	}

	disableFallback := os.Getenv("ANY_REGISTRY_DISABLE")
	if disableFallback == "true" {
		return "", "", err
	}

	username, password, fallbackErr := e.getFor("ANY_REGISTRY")
	if fallbackErr != nil {
		return "", "", fmt.Errorf("failed to get fallback credentials for %s (set ANY_REGISTRY_DISABLE=true to disable fallback): %w", serverURL, fallbackErr)
	}
	return username, password, nil

}

// List is not supported and will always error.
func (*Env) List() (map[string]string, error) {
	return nil, errors.New("list is not supported")
}
