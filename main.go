/*
Copyright 2021 Matt Moore
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"knative.dev/pkg/signals"

	"github.com/mattmoor/oidc-magic/pkg/providers"

	// These are the registered plugins
	_ "github.com/mattmoor/oidc-magic/pkg/providers/github"
	_ "github.com/mattmoor/oidc-magic/pkg/providers/google"
)

func main() {
	ctx := signals.NewContext()

	if providers.Enabled(ctx) {
		tok, err := providers.Provide(ctx, "sigstore")
		if err != nil {
			log.Fatalf("Provide() = %v", err)
		}
		log.Printf("Providers are enabled: %v", tok)
	} else {
		log.Print("Providers are not enabled :(")
	}
}
