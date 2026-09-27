// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vaccount

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
	_ datasource.DataSource              = (*contactsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*contactsDataSource)(nil)
)

// ContactsDataSource returns the hitechcloud_contacts data source.
func ContactsDataSource() datasource.DataSource {
	return &contactsDataSource{}
}

type contactsDataSource struct {
	client *client.Client
}

type contactsModel struct {
	ID       types.String  `tfsdk:"id"`
	Contacts []contactItem `tfsdk:"contacts"`
}

type contactItem struct {
	ID          types.String `tfsdk:"id"`
	Email       types.String `tfsdk:"email"`
	Type        types.String `tfsdk:"type"`
	FirstName   types.String `tfsdk:"first_name"`
	LastName    types.String `tfsdk:"last_name"`
	CompanyName types.String `tfsdk:"company_name"`
	PhoneNumber types.String `tfsdk:"phone_number"`
	Country     types.String `tfsdk:"country"`
	State       types.String `tfsdk:"state"`
	City        types.String `tfsdk:"city"`
	Address     types.String `tfsdk:"address"`
	Postcode    types.String `tfsdk:"postcode"`
}

func (d *contactsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contacts"
}

func (d *contactsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the sub-account contacts of the HiTechCloud account " +
			"(`GET /api/contact`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (`contacts`).",
			},
			"contacts": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Contacts of the account.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":           schema.StringAttribute{Computed: true, MarkdownDescription: "Contact identifier."},
						"email":        schema.StringAttribute{Computed: true, MarkdownDescription: "Contact email."},
						"type":         schema.StringAttribute{Computed: true, MarkdownDescription: "Contact type."},
						"first_name":   schema.StringAttribute{Computed: true, MarkdownDescription: "First name."},
						"last_name":    schema.StringAttribute{Computed: true, MarkdownDescription: "Last name."},
						"company_name": schema.StringAttribute{Computed: true, MarkdownDescription: "Company name."},
						"phone_number": schema.StringAttribute{Computed: true, MarkdownDescription: "Phone number."},
						"country":      schema.StringAttribute{Computed: true, MarkdownDescription: "Country."},
						"state":        schema.StringAttribute{Computed: true, MarkdownDescription: "State / province."},
						"city":         schema.StringAttribute{Computed: true, MarkdownDescription: "City."},
						"address":      schema.StringAttribute{Computed: true, MarkdownDescription: "Street address."},
						"postcode":     schema.StringAttribute{Computed: true, MarkdownDescription: "Postal code."},
					},
				},
			},
		},
	}
}

func (d *contactsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *contactsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	contacts, err := d.client.ListContacts(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing contacts", err.Error())
		return
	}

	items := make([]contactItem, 0, len(contacts))
	for _, c := range contacts {
		items = append(items, contactItem{
			ID:          common.StringOrNull(c.ID),
			Email:       common.StringOrNull(c.Email),
			Type:        common.StringOrNull(c.Type),
			FirstName:   common.StringOrNull(c.FirstName),
			LastName:    common.StringOrNull(c.LastName),
			CompanyName: common.StringOrNull(c.CompanyName),
			PhoneNumber: common.StringOrNull(c.PhoneNumber),
			Country:     common.StringOrNull(c.Country),
			State:       common.StringOrNull(c.State),
			City:        common.StringOrNull(c.City),
			Address:     common.StringOrNull(c.Address),
			Postcode:    common.StringOrNull(c.Postcode),
		})
	}

	data := contactsModel{
		ID:       types.StringValue("contacts"),
		Contacts: items,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
