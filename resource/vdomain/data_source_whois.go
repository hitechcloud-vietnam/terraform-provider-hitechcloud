// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vdomain

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/client"
)

var (
	_ datasource.DataSource              = (*whoisDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*whoisDataSource)(nil)
)

// WhoisDataSource returns the hitechcloud_whois data source.
func WhoisDataSource() datasource.DataSource {
	return &whoisDataSource{}
}

type whoisDataSource struct {
	client *client.Client
}

type whoisModel struct {
	ID         types.String `tfsdk:"id"`
	Domain     types.String `tfsdk:"domain"`
	Attributes types.Map    `tfsdk:"attributes"`
}

func (d *whoisDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_whois"
}

func (d *whoisDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Performs a WHOIS lookup for a domain (`GET /api/whois/{domain}`). " +
			"The raw WHOIS fields are exposed as the `attributes` map.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (same as `domain`).",
			},
			"domain": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Domain name to look up.",
			},
			"attributes": schema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "WHOIS attributes of the domain.",
			},
		},
	}
}

func (d *whoisDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *whoisDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data whoisModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := data.Domain.ValueString()
	attrs, err := d.client.GetWhois(ctx, domain)
	if err != nil {
		resp.Diagnostics.AddError("Error performing WHOIS lookup", err.Error())
		return
	}

	attrMap, diags := types.MapValueFrom(ctx, types.StringType, attrs)
	resp.Diagnostics.Append(diags...)

	data.ID = types.StringValue(domain)
	data.Attributes = attrMap
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
