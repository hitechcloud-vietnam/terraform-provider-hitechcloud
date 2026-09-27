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
	_ datasource.DataSource              = (*domainTLDsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*domainTLDsDataSource)(nil)
)

// DomainTLDsDataSource returns the hitechcloud_domain_tlds data source.
func DomainTLDsDataSource() datasource.DataSource {
	return &domainTLDsDataSource{}
}

type domainTLDsDataSource struct {
	client *client.Client
}

type domainTLDsModel struct {
	ID   types.String `tfsdk:"id"`
	TLDs []tldItem    `tfsdk:"tlds"`
}

type tldItem struct {
	ID       types.String `tfsdk:"id"`
	TLD      types.String `tfsdk:"tld"`
	MinYears types.Int64  `tfsdk:"min_years"`
	MaxYears types.Int64  `tfsdk:"max_years"`
}

func (d *domainTLDsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_tlds"
}

func (d *domainTLDsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the TLDs available for registration and transfer " +
			"(`GET /api/domain/order`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (`domain_tlds`).",
			},
			"tlds": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Available TLDs.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":        schema.StringAttribute{Computed: true, MarkdownDescription: "TLD record identifier."},
						"tld":       schema.StringAttribute{Computed: true, MarkdownDescription: "TLD extension, e.g. `.vn`."},
						"min_years": schema.Int64Attribute{Computed: true, MarkdownDescription: "Minimum registration period in years."},
						"max_years": schema.Int64Attribute{Computed: true, MarkdownDescription: "Maximum registration period in years."},
					},
				},
			},
		},
	}
}

func (d *domainTLDsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *domainTLDsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tlds, err := d.client.ListDomainTLDs(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing domain TLDs", err.Error())
		return
	}

	items := make([]tldItem, 0, len(tlds))
	for _, t := range tlds {
		items = append(items, tldItem{
			ID:       common.StringOrNull(t.ID),
			TLD:      common.StringOrNull(t.TLD),
			MinYears: common.Int64OrNull(t.MinYears),
			MaxYears: common.Int64OrNull(t.MaxYears),
		})
	}

	data := domainTLDsModel{
		ID:   types.StringValue("domain_tlds"),
		TLDs: items,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
