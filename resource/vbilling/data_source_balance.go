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
	_ datasource.DataSource              = (*balanceDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*balanceDataSource)(nil)
)

// BalanceDataSource returns the hitechcloud_balance data source.
func BalanceDataSource() datasource.DataSource {
	return &balanceDataSource{}
}

type balanceDataSource struct {
	client *client.Client
}

type balanceModel struct {
	ID       types.String `tfsdk:"id"`
	Credit   types.String `tfsdk:"credit"`
	Balance  types.String `tfsdk:"balance"`
	Currency types.String `tfsdk:"currency"`
}

func (d *balanceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_balance"
}

func (d *balanceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves the account balance of the HiTechCloud account " +
			"(`GET /api/balance`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (`balance`).",
			},
			"credit":   schema.StringAttribute{Computed: true, MarkdownDescription: "Available credit."},
			"balance":  schema.StringAttribute{Computed: true, MarkdownDescription: "Account balance."},
			"currency": schema.StringAttribute{Computed: true, MarkdownDescription: "Currency."},
		},
	}
}

func (d *balanceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *balanceDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	bal, err := d.client.GetBalance(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading balance", err.Error())
		return
	}

	data := balanceModel{
		ID:       types.StringValue("balance"),
		Credit:   common.StringOrNull(bal.Credit),
		Balance:  common.StringOrNull(bal.Balance),
		Currency: common.StringOrNull(bal.Currency),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
