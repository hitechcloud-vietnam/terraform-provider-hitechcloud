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
	_ datasource.DataSource              = (*domainAvailabilityDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*domainAvailabilityDataSource)(nil)
)

// DomainAvailabilityDataSource returns the hitechcloud_domain_availability data source.
func DomainAvailabilityDataSource() datasource.DataSource {
	return &domainAvailabilityDataSource{}
}

type domainAvailabilityDataSource struct {
	client *client.Client
}

type domainAvailabilityModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Available types.Bool   `tfsdk:"available"`
	Price     types.String `tfsdk:"price"`
}

func (d *domainAvailabilityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_availability"
}

func (d *domainAvailabilityDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Checks whether a domain name is available for registration " +
			"(`POST /api/domain/lookup`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (same as `name`).",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Domain name to check.",
			},
			"available": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the domain is available for registration.",
			},
			"price": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Registration price, when returned by the API.",
			},
		},
	}
}

func (d *domainAvailabilityDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *domainAvailabilityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data domainAvailabilityModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	available, price, err := d.client.CheckDomainAvailability(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError("Error checking domain availability", err.Error())
		return
	}

	data.ID = types.StringValue(name)
	data.Available = types.BoolValue(available)
	data.Price = types.StringValue(price)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
