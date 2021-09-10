/*
Copyright 2021 Matt Moore
SPDX-License-Identifier: Apache-2.0
*/

package google

import (
	"context"
	"io/ioutil"
	"strings"

	"google.golang.org/api/idtoken"
	"knative.dev/pkg/logging"

	"github.com/mattmoor/oidc-magic/pkg/providers"
)

func init() {
	providers.Register("google-workload-identity", &googleWorkloadIdentity{})
}

type googleWorkloadIdentity struct{}

var _ providers.Interface = (*googleWorkloadIdentity)(nil)

// gceProductNameFile is the product file path that contains the cloud service name.
// This is a variable instead of a const to enable testing.
var gceProductNameFile = "/sys/class/dmi/id/product_name"

// Enabled implements providers.Interface
// This is based on k8s.io/kubernetes/pkg/credentialprovider/gcp
func (gwi *googleWorkloadIdentity) Enabled(ctx context.Context) bool {
	data, err := ioutil.ReadFile(gceProductNameFile)
	if err != nil {
		logging.FromContext(ctx).Debugf("Error while reading product_name: %v", err)
		return false
	}
	name := strings.TrimSpace(string(data))
	if name == "Google" || name == "Google Compute Engine" {
		// Just because we're on Google, does not mean workload identity is available.
		// TODO(mattmoor): do something better than this.
		_, err := gwi.Provide(ctx, "garbage")
		return err == nil
	}
	return false
}

// Provide implements providers.Interface
func (gwi *googleWorkloadIdentity) Provide(ctx context.Context, audience string) (string, error) {
	ts, err := idtoken.NewTokenSource(ctx, audience)
	if err != nil {
		return "", err
	}
	tok, err := ts.Token()
	if err != nil {
		return "", err
	}
	return tok.AccessToken, nil
}
