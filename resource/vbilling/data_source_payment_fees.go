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
	_ datasource.DataSource              = (*paymentFeesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*paymentFeesDataSource)(nil)
)

// PaymentFeesDataSource returns the hitechcloud_payment_fees data source.
func PaymentFeesDataSource() datasource.DataSource {
	return &paymentFeesDataSource{}
}

type paymentFeesDataSource struct {
	client *client.Client
}

type paymentFeesModel struct {
	ID   types.String     `tfsdk:"id"`
	Fees []paymentFeeItem `tfsdk:"fees"`
}

type paymentFeeItem struct {
	Module types.String `tfsdk:"module"`
	Name   types.String `tfsdk:"name"`
	Fee    types.String `tfsdk:"fee"`
}

func (d *paymentFeesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_payment_fees"
}

func (d *paymentFeesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the fees of the payment methods available to the " +
			"HiTechCloud account (`GET /api/payment/fees`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (`payment_fees`).",
			},
			"fees": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Payment method fees.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"module": schema.StringAttribute{Computed: true, MarkdownDescription: "Payment module / gateway key."},
						"name":   schema.StringAttribute{Computed: true, MarkdownDescription: "Display name."},
						"fee":    schema.StringAttribute{Computed: true, MarkdownDescription: "Fee (percentage or amount)."},
					},
				},
			},
		},
	}
}

func (d *paymentFeesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *paymentFeesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	fees, err := d.client.ListPaymentFees(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing payment fees", err.Error())
		return
	}

	items := make([]paymentFeeItem, 0, len(fees))
	for _, f := range fees {
		items = append(items, paymentFeeItem{
			Module: common.StringOrNull(f.Module),
			Name:   common.StringOrNull(f.Name),
			Fee:    common.StringOrNull(f.Fee),
		})
	}

	data := paymentFeesModel{
		ID:   types.StringValue("payment_fees"),
		Fees: items,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
