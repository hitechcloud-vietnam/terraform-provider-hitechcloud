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
	_ datasource.DataSource              = (*domainsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*domainsDataSource)(nil)
)

// DomainsDataSource returns the hitechcloud_domains data source.
func DomainsDataSource() datasource.DataSource {
	return &domainsDataSource{}
}

type domainsDataSource struct {
	client *client.Client
}

type domainsModel struct {
	ID      types.String `tfsdk:"id"`
	Domains []domainItem `tfsdk:"domains"`
}

// domainItem is the shared nested object for domain listings.
type domainItem struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Status        types.String `tfsdk:"status"`
	Registrar     types.String `tfsdk:"registrar"`
	RegDate       types.String `tfsdk:"reg_date"`
	Expires       types.String `tfsdk:"expires"`
	Autorenew     types.Bool   `tfsdk:"autorenew"`
	RegistrarLock types.Bool   `tfsdk:"registrar_lock"`
	IDProtection  types.Bool   `tfsdk:"id_protection"`
	Nameservers   []string     `tfsdk:"nameservers"`
	EPPCode       types.String `tfsdk:"epp_code"`
}

func domainItemAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true, MarkdownDescription: "Domain identifier."},
		"name":           schema.StringAttribute{Computed: true, MarkdownDescription: "Domain name."},
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
	}
}

func toDomainItem(d client.Domain) domainItem {
	return domainItem{
		ID:            common.StringOrNull(d.ID),
		Name:          common.StringOrNull(d.Name),
		Status:        common.StringOrNull(d.Status),
		Registrar:     common.StringOrNull(d.Registrar),
		RegDate:       common.StringOrNull(d.RegDate),
		Expires:       common.StringOrNull(d.Expires),
		Autorenew:     types.BoolValue(d.Autorenew),
		RegistrarLock: types.BoolValue(d.RegistrarLock),
		IDProtection:  types.BoolValue(d.IDProtection),
		Nameservers:   d.Nameservers,
		EPPCode:       common.StringOrNull(d.EPPCode),
	}
}

func (d *domainsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domains"
}

func (d *domainsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the registered domains of the HiTechCloud account (`GET /api/domain`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (`domains`).",
			},
			"domains": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Registered domains.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: domainItemAttributes(),
				},
			},
		},
	}
}

func (d *domainsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *domainsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	domains, err := d.client.ListDomains(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing domains", err.Error())
		return
	}

	items := make([]domainItem, 0, len(domains))
	for _, dom := range domains {
		items = append(items, toDomainItem(dom))
	}

	data := domainsModel{
		ID:      types.StringValue("domains"),
		Domains: items,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
