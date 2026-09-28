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
	_ datasource.DataSource              = (*domainContactDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*domainContactDataSource)(nil)
)

// DomainContactDataSource returns the hitechcloud_domain_contact data source.
func DomainContactDataSource() datasource.DataSource {
	return &domainContactDataSource{}
}

type domainContactDataSource struct {
	client *client.Client
}

type domainContactModel struct {
	ID          types.String `tfsdk:"id"`
	DomainID    types.String `tfsdk:"domain_id"`
	Registrant  types.String `tfsdk:"registrant_contact_id"`
	Admin       types.String `tfsdk:"admin_contact_id"`
	Tech        types.String `tfsdk:"tech_contact_id"`
	Billing     types.String `tfsdk:"billing_contact_id"`
	EPPCode     types.String `tfsdk:"epp_code"`
	Locked      types.Bool   `tfsdk:"registrar_lock"`
	IDProtected types.Bool   `tfsdk:"id_protection"`
}

func (d *domainContactDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_contact"
}

func (d *domainContactDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads contact and transfer information of a registered domain " +
			"(`GET /api/domain/{id}/contact`, `GET /api/domain/{id}/epp`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The domain ID (same as `domain_id`).",
			},
			"domain_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the registered domain.",
			},
			"registrant_contact_id": schema.StringAttribute{Computed: true, MarkdownDescription: "Registrant contact ID."},
			"admin_contact_id":      schema.StringAttribute{Computed: true, MarkdownDescription: "Admin contact ID."},
			"tech_contact_id":       schema.StringAttribute{Computed: true, MarkdownDescription: "Tech contact ID."},
			"billing_contact_id":    schema.StringAttribute{Computed: true, MarkdownDescription: "Billing contact ID."},
			"epp_code":              schema.StringAttribute{Computed: true, Sensitive: true, MarkdownDescription: "EPP transfer code of the domain."},
			"registrar_lock":        schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the domain is registrar-locked."},
			"id_protection":         schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether WHOIS ID protection is enabled."},
		},
	}
}

func (d *domainContactDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *domainContactDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data domainContactModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainID := data.DomainID.ValueString()
	data.ID = types.StringValue(domainID)

	contact, err := d.client.GetDomainContactInfo(ctx, domainID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading domain contact info", err.Error())
		return
	}
	data.Registrant = common.StringOrNull(client.FirstString(contact, "registrant", "registrant_id", "registrant_contact_id"))
	data.Admin = common.StringOrNull(client.FirstString(contact, "admin", "admin_id", "admin_contact_id"))
	data.Tech = common.StringOrNull(client.FirstString(contact, "tech", "tech_id", "tech_contact_id"))
	data.Billing = common.StringOrNull(client.FirstString(contact, "billing", "billing_id", "billing_contact_id"))
	data.Locked = types.BoolValue(client.FirstBool(contact, "registrar_lock", "registrarlock", "locked"))
	data.IDProtected = types.BoolValue(client.FirstBool(contact, "id_protection", "idprotection", "id_protect"))

	if epp, err := d.client.GetDomainEPPCode(ctx, domainID); err == nil {
		data.EPPCode = common.StringOrNull(client.FirstString(epp, "epp", "code", "epp_code"))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
