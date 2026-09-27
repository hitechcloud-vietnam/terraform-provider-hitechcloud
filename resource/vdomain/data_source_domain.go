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
	_ datasource.DataSource              = (*domainDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*domainDataSource)(nil)
)

// DomainDataSource returns the hitechcloud_domain data source.
func DomainDataSource() datasource.DataSource {
	return &domainDataSource{}
}

type domainDataSource struct {
	client *client.Client
}

type domainModel struct {
	ID            types.String `tfsdk:"id"`
	DomainID      types.String `tfsdk:"domain_id"`
	Name          types.String `tfsdk:"name"`
	Status        types.String `tfsdk:"status"`
	Registrar     types.String `tfsdk:"registrar"`
	RegDate       types.String `tfsdk:"reg_date"`
	Expires       types.String `tfsdk:"expires"`
	Autorenew     types.Bool   `tfsdk:"autorenew"`
	RegistrarLock types.Bool   `tfsdk:"registrar_lock"`
	IDProtection  types.Bool   `tfsdk:"id_protection"`
	Nameservers   types.Set    `tfsdk:"nameservers"`
	EPPCode       types.String `tfsdk:"epp_code"`
}

func (d *domainDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (d *domainDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads a single registered domain, either by `domain_id` " +
			"(`GET /api/domain/{id}`) or by `name` (`GET /api/domain/name/{name}`). " +
			"Exactly one of the two must be set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (the domain id).",
			},
			"domain_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Domain identifier to look up (mutually exclusive with `name`).",
			},
			"name":           schema.StringAttribute{Optional: true, MarkdownDescription: "Domain name to look up (mutually exclusive with `domain_id`)."},
			"status":         schema.StringAttribute{Computed: true, MarkdownDescription: "Domain status."},
			"registrar":      schema.StringAttribute{Computed: true, MarkdownDescription: "Registrar."},
			"reg_date":       schema.StringAttribute{Computed: true, MarkdownDescription: "Registration date."},
			"expires":        schema.StringAttribute{Computed: true, MarkdownDescription: "Expiry date."},
			"autorenew":      schema.BoolAttribute{Computed: true, MarkdownDescription: "Auto-renew enabled."},
			"registrar_lock": schema.BoolAttribute{Computed: true, MarkdownDescription: "Registrar lock enabled."},
			"id_protection":  schema.BoolAttribute{Computed: true, MarkdownDescription: "ID protection enabled."},
			"nameservers": schema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Nameservers.",
			},
			"epp_code": schema.StringAttribute{Computed: true, MarkdownDescription: "EPP / transfer secret code."},
		},
	}
}

func (d *domainDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *domainDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data domainModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainID := data.DomainID.ValueString()
	name := data.Name.ValueString()

	if (domainID == "") == (name == "") {
		resp.Diagnostics.AddError(
			"Invalid combination of lookup fields",
			"Exactly one of `domain_id` or `name` must be set.",
		)
		return
	}

	var (
		dom *client.Domain
		err error
	)
	if domainID != "" {
		dom, err = d.client.GetDomain(ctx, domainID)
	} else {
		dom, err = d.client.GetDomainByName(ctx, name)
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading domain", err.Error())
		return
	}

	item := toDomainItem(*dom)
	data.ID = item.ID
	data.DomainID = item.ID
	data.Name = item.Name
	data.Status = item.Status
	data.Registrar = item.Registrar
	data.RegDate = item.RegDate
	data.Expires = item.Expires
	data.Autorenew = item.Autorenew
	data.RegistrarLock = item.RegistrarLock
	data.IDProtection = item.IDProtection
	data.Nameservers = common.SetStrings(item.Nameservers)
	data.EPPCode = item.EPPCode

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
