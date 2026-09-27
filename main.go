// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/provider"
)

// version and commit are set by the release pipeline (GoReleaser) via ldflags.
var (
	version = "dev"
	commit  = "none"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	if debug {
		log.Printf("terraform-provider-hitechcloud version=%s commit=%s", version, commit)
	}

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/hitechcloud-vietnam/hitechcloud",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err)
	}
}
