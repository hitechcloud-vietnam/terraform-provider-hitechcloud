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
	_ datasource.DataSource              = (*accountDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*accountDataSource)(nil)
)

// AccountDataSource returns the hitechcloud_account data source.
func AccountDataSource() datasource.DataSource {
	return &accountDataSource{}
}

type accountDataSource struct {
	client *client.Client
}

type accountModel struct {
	ID            types.String `tfsdk:"id"`
	Email         types.String `tfsdk:"email"`
	Type          types.String `tfsdk:"type"`
	FirstName     types.String `tfsdk:"first_name"`
	LastName      types.String `tfsdk:"last_name"`
	CompanyName   types.String `tfsdk:"company_name"`
	TaxID         types.String `tfsdk:"tax_id"`
	Gender        types.String `tfsdk:"gender"`
	NationalID    types.String `tfsdk:"national_id"`
	Birthday      types.String `tfsdk:"birthday"`
	PhoneNumber   types.String `tfsdk:"phone_number"`
	Country       types.String `tfsdk:"country"`
	State         types.String `tfsdk:"state"`
	City          types.String `tfsdk:"city"`
	Address       types.String `tfsdk:"address"`
	Postcode      types.String `tfsdk:"postcode"`
	Currency      types.String `tfsdk:"currency"`
	BankName      types.String `tfsdk:"bank_name"`
	BankAccount   types.String `tfsdk:"bank_account"`
	BankAccountNm types.String `tfsdk:"bank_account_name"`
}

func (d *accountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

func (d *accountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads the HiTechCloud account that owns the API token " +
			"(`GET /api/details`).",
		Attributes: map[string]schema.Attribute{
			"id":                schema.StringAttribute{Computed: true, MarkdownDescription: "Account identifier."},
			"email":             schema.StringAttribute{Computed: true, MarkdownDescription: "Account email address."},
			"type":              schema.StringAttribute{Computed: true, MarkdownDescription: "Account type."},
			"first_name":        schema.StringAttribute{Computed: true, MarkdownDescription: "First name."},
			"last_name":         schema.StringAttribute{Computed: true, MarkdownDescription: "Last name."},
			"company_name":      schema.StringAttribute{Computed: true, MarkdownDescription: "Company name."},
			"tax_id":            schema.StringAttribute{Computed: true, MarkdownDescription: "Tax code."},
			"gender":            schema.StringAttribute{Computed: true, MarkdownDescription: "Gender."},
			"national_id":       schema.StringAttribute{Computed: true, MarkdownDescription: "National ID (CCCD)."},
			"birthday":          schema.StringAttribute{Computed: true, MarkdownDescription: "Birthday."},
			"phone_number":      schema.StringAttribute{Computed: true, MarkdownDescription: "Phone number."},
			"country":           schema.StringAttribute{Computed: true, MarkdownDescription: "Country."},
			"state":             schema.StringAttribute{Computed: true, MarkdownDescription: "State / province."},
			"city":              schema.StringAttribute{Computed: true, MarkdownDescription: "City."},
			"address":           schema.StringAttribute{Computed: true, MarkdownDescription: "Street address."},
			"postcode":          schema.StringAttribute{Computed: true, MarkdownDescription: "Postal code."},
			"currency":          schema.StringAttribute{Computed: true, MarkdownDescription: "Account currency."},
			"bank_name":         schema.StringAttribute{Computed: true, MarkdownDescription: "Bank name."},
			"bank_account":      schema.StringAttribute{Computed: true, MarkdownDescription: "Bank account number."},
			"bank_account_name": schema.StringAttribute{Computed: true, MarkdownDescription: "Bank account holder name."},
		},
	}
}

func (d *accountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *accountDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	acc, err := d.client.GetAccount(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading account details", err.Error())
		return
	}

	data := accountModel{
		ID:            common.StringOrNull(acc.ID),
		Email:         common.StringOrNull(acc.Email),
		Type:          common.StringOrNull(acc.Type),
		FirstName:     common.StringOrNull(acc.FirstName),
		LastName:      common.StringOrNull(acc.LastName),
		CompanyName:   common.StringOrNull(acc.CompanyName),
		TaxID:         common.StringOrNull(acc.TaxID),
		Gender:        common.StringOrNull(acc.Gender),
		NationalID:    common.StringOrNull(acc.NationalID),
		Birthday:      common.StringOrNull(acc.Birthday),
		PhoneNumber:   common.StringOrNull(acc.PhoneNumber),
		Country:       common.StringOrNull(acc.Country),
		State:         common.StringOrNull(acc.State),
		City:          common.StringOrNull(acc.City),
		Address:       common.StringOrNull(acc.Address),
		Postcode:      common.StringOrNull(acc.Postcode),
		Currency:      common.StringOrNull(acc.Currency),
		BankName:      common.StringOrNull(acc.BankName),
		BankAccount:   common.StringOrNull(acc.BankAccount),
		BankAccountNm: common.StringOrNull(acc.BankAccountNm),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
