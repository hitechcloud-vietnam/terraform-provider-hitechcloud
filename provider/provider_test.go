// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package provider_test

import (
	"context"
	"testing"

	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"

	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/provider"
)

func TestProviderMetadata(t *testing.T) {
	p := provider.New("test")()
	var resp fwprovider.MetadataResponse
	p.Metadata(context.Background(), fwprovider.MetadataRequest{}, &resp)

	if resp.TypeName != "hitechcloud" {
		t.Fatalf("TypeName = %q, want hitechcloud", resp.TypeName)
	}
	if resp.Version != "test" {
		t.Fatalf("Version = %q, want test", resp.Version)
	}
}

func TestProviderSchemaHasConfigAttributes(t *testing.T) {
	p := provider.New("test")()
	var resp fwprovider.SchemaResponse
	p.Schema(context.Background(), fwprovider.SchemaRequest{}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() diagnostics: %v", resp.Diagnostics)
	}
	for _, attr := range []string{"token", "endpoint", "request_timeout"} {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Fatalf("schema is missing attribute %q", attr)
		}
	}
}

func TestProviderRegistersResourcesAndDataSources(t *testing.T) {
	p := provider.New("test")()

	resources := p.Resources(context.Background())
	if len(resources) != 19 {
		t.Fatalf("Resources() count = %d, want 19", len(resources))
	}
	dataSources := p.DataSources(context.Background())
	if len(dataSources) != 30 {
		t.Fatalf("DataSources() count = %d, want 30", len(dataSources))
	}

	// Every registered type must be instantiable.
	for i, f := range resources {
		if f() == nil {
			t.Fatalf("Resources()[%d] returned nil", i)
		}
	}
	for i, f := range dataSources {
		if f() == nil {
			t.Fatalf("DataSources()[%d] returned nil", i)
		}
	}
}
