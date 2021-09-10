/*
Copyright 2021 Matt Moore
SPDX-License-Identifier: Apache-2.0
*/

package github

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/mattmoor/oidc-magic/pkg/providers"
)

func init() {
	providers.Register("github-actions", &githubActions{})
}

type githubActions struct{}

var _ providers.Interface = (*githubActions)(nil)

const (
	RequestTokenEnvKey = "ACTIONS_ID_TOKEN_REQUEST_TOKEN"
	RequestURLEnvKey   = "ACTIONS_ID_TOKEN_REQUEST_URL"
)

// Enabled implements providers.Interface
func (ga *githubActions) Enabled(ctx context.Context) bool {
	if os.Getenv(RequestTokenEnvKey) == "" {
		return false
	}
	if os.Getenv(RequestURLEnvKey) == "" {
		return false
	}
	return true
}

// Provide implements providers.Interface
func (ga *githubActions) Provide(ctx context.Context, audience string) (string, error) {
	url := os.Getenv(RequestURLEnvKey) + "&audience=" + audience

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "bearer "+os.Getenv(RequestTokenEnvKey))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var payload struct {
		Value string `json:"value"`
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&payload); err != nil {
		return "", err
	}
	return payload.Value, nil
}
