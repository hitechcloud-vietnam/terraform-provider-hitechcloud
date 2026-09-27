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
	_ datasource.DataSource              = (*domainDNSTypesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*domainDNSTypesDataSource)(nil)
)

// DomainDNSTypesDataSource returns the hitechcloud_domain_dns_types data source.
func DomainDNSTypesDataSource() datasource.DataSource {
	return &domainDNSTypesDataSource{}
}

type domainDNSTypesDataSource struct {
	client *client.Client
}

type domainDNSTypesModel struct {
	ID       types.String `tfsdk:"id"`
	DomainID types.String `tfsdk:"domain_id"`
	Types    types.Set    `tfsdk:"types"`
}

func (d *domainDNSTypesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_dns_types"
}

func (d *domainDNSTypesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the DNS record types supported by a registered " +
			"HiTechCloud domain (`GET /api/domain/{domain_id}/dns/types`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (same as `domain_id`).",
			},
			"domain_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the registered domain.",
			},
			"types": schema.SetAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Supported DNS record types.",
			},
		},
	}
}

func (d *domainDNSTypesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *domainDNSTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data domainDNSTypesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainID := data.DomainID.ValueString()
	typesList, err := d.client.ListDomainDNSTypes(ctx, domainID)
	if err != nil {
		resp.Diagnostics.AddError("Error listing domain DNS types", err.Error())
		return
	}

	set, diags := types.SetValueFrom(ctx, types.StringType, typesList)
	resp.Diagnostics.Append(diags...)

	data.ID = types.StringValue(domainID)
	data.Types = set
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
