// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vbilling

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
	_ datasource.DataSource              = (*certificatesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*certificatesDataSource)(nil)
)

// CertificatesDataSource returns the hitechcloud_certificates data source.
func CertificatesDataSource() datasource.DataSource {
	return &certificatesDataSource{}
}

type certificatesDataSource struct {
	client *client.Client
}

type certificatesModel struct {
	ID           types.String      `tfsdk:"id"`
	Certificates []certificateItem `tfsdk:"certificates"`
}

type certificateItem struct {
	ID      types.String `tfsdk:"id"`
	Domain  types.String `tfsdk:"domain"`
	Status  types.String `tfsdk:"status"`
	Type    types.String `tfsdk:"type"`
	Expires types.String `tfsdk:"expires"`
}

func (d *certificatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificates"
}

func (d *certificatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the SSL certificates of the HiTechCloud account " +
			"(`GET /api/certificate`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (`certificates`).",
			},
			"certificates": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "SSL certificates.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":      schema.StringAttribute{Computed: true, MarkdownDescription: "Certificate identifier."},
						"domain":  schema.StringAttribute{Computed: true, MarkdownDescription: "Certificate domain."},
						"status":  schema.StringAttribute{Computed: true, MarkdownDescription: "Certificate status."},
						"type":    schema.StringAttribute{Computed: true, MarkdownDescription: "Certificate product / type."},
						"expires": schema.StringAttribute{Computed: true, MarkdownDescription: "Expiry date."},
					},
				},
			},
		},
	}
}

func (d *certificatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *certificatesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	certs, err := d.client.ListCertificates(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing certificates", err.Error())
		return
	}

	items := make([]certificateItem, 0, len(certs))
	for _, c := range certs {
		items = append(items, certificateItem{
			ID:      common.StringOrNull(c.ID),
			Domain:  common.StringOrNull(c.Domain),
			Status:  common.StringOrNull(c.Status),
			Type:    common.StringOrNull(c.Type),
			Expires: common.StringOrNull(c.Expires),
		})
	}

	data := certificatesModel{
		ID:           types.StringValue("certificates"),
		Certificates: items,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
