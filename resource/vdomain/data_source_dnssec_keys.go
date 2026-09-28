// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vdomain

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/client"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/common"
)

var (
	_ datasource.DataSource              = (*dnssecKeysDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*dnssecKeysDataSource)(nil)
)

// DNSSECKeysDataSource returns the hitechcloud_dnssec_keys data source.
func DNSSECKeysDataSource() datasource.DataSource {
	return &dnssecKeysDataSource{}
}

type dnssecKeysDataSource struct {
	client *client.Client
}

type dnssecKeysModel struct {
	ID             types.String     `tfsdk:"id"`
	DomainID       types.String     `tfsdk:"domain_id"`
	Keys           []dnssecKeyEntry `tfsdk:"keys"`
	AvailableFlags []types.String   `tfsdk:"available_flags"`
}

type dnssecKeyEntry struct {
	KeyTag     types.String `tfsdk:"key_tag"`
	Algorithm  types.String `tfsdk:"algorithm"`
	DigestType types.String `tfsdk:"digest_type"`
	Digest     types.String `tfsdk:"digest"`
	Flags      types.String `tfsdk:"flags"`
	Protocol   types.String `tfsdk:"protocol"`
	PublicKey  types.String `tfsdk:"public_key"`
}

func (d *dnssecKeysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dnssec_keys"
}

func (d *dnssecKeysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the DNSSEC keys (DS records) and available flags of a " +
			"registered domain (`GET /api/domain/{id}/dnssec`, `GET /api/domain/{id}/dnssec/flags`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The domain ID (same as `domain_id`).",
			},
			"domain_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the registered domain.",
			},
			"keys": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "DNSSEC keys configured for the domain.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key_tag":     schema.StringAttribute{Computed: true, MarkdownDescription: "DS key tag."},
						"algorithm":   schema.StringAttribute{Computed: true, MarkdownDescription: "DNSSEC algorithm number."},
						"digest_type": schema.StringAttribute{Computed: true, MarkdownDescription: "DS digest type."},
						"digest":      schema.StringAttribute{Computed: true, MarkdownDescription: "Hex-encoded DS digest."},
						"flags":       schema.StringAttribute{Computed: true, MarkdownDescription: "DNSKEY flags."},
						"protocol":    schema.StringAttribute{Computed: true, MarkdownDescription: "DNSKEY protocol number."},
						"public_key":  schema.StringAttribute{Computed: true, MarkdownDescription: "Base64-encoded public key material."},
					},
				},
			},
			"available_flags": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Flags the API accepts for DNSSEC keys.",
			},
		},
	}
}

func (d *dnssecKeysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	cli, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got %T", req.ProviderData),
		)
		return
	}
	d.client = cli
}

func (d *dnssecKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dnssecKeysModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainID := data.DomainID.ValueString()
	data.ID = types.StringValue(domainID)

	keys, err := d.client.ListDNSSECKeys(ctx, domainID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading DNSSEC keys", err.Error())
		return
	}
	data.Keys = make([]dnssecKeyEntry, 0, len(keys))
	for _, item := range keys {
		m, ok := client.AsMap(item)
		if !ok {
			continue
		}
		data.Keys = append(data.Keys, dnssecKeyEntry{
			KeyTag:     common.StringOrNull(client.FirstString(m, "key_tag", "keytag", "tag")),
			Algorithm:  common.StringOrNull(client.FirstString(m, "algorithm", "alg")),
			DigestType: common.StringOrNull(client.FirstString(m, "digest_type", "digesttype")),
			Digest:     common.StringOrNull(client.FirstString(m, "digest")),
			Flags:      common.StringOrNull(client.FirstString(m, "flags", "flag")),
			Protocol:   common.StringOrNull(client.FirstString(m, "protocol", "proto")),
			PublicKey:  common.StringOrNull(client.FirstString(m, "public_key", "pubkey", "key")),
		})
	}

	flagValues := []string{}
	if flags, err := d.client.ListDNSSECFlags(ctx, domainID); err == nil {
		for _, f := range flags {
			if s := client.FirstString(f, "flag", "value", "name"); s != "" {
				flagValues = append(flagValues, s)
			}
		}
	}
	data.AvailableFlags = make([]types.String, 0, len(flagValues))
	for _, f := range flagValues {
		data.AvailableFlags = append(data.AvailableFlags, types.StringValue(f))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
