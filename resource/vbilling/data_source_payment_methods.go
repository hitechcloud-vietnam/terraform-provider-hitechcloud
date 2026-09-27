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
	_ datasource.DataSource              = (*paymentMethodsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*paymentMethodsDataSource)(nil)
)

// PaymentMethodsDataSource returns the hitechcloud_payment_methods data source.
func PaymentMethodsDataSource() datasource.DataSource {
	return &paymentMethodsDataSource{}
}

type paymentMethodsDataSource struct {
	client *client.Client
}

type paymentMethodsModel struct {
	ID      types.String        `tfsdk:"id"`
	Methods []paymentMethodItem `tfsdk:"methods"`
}

type paymentMethodItem struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *paymentMethodsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_payment_methods"
}

func (d *paymentMethodsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the payment methods available to the HiTechCloud " +
			"account (`GET /api/payment`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (`payment_methods`).",
			},
			"methods": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Available payment methods.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true, MarkdownDescription: "Payment module / gateway key."},
						"name": schema.StringAttribute{Computed: true, MarkdownDescription: "Display name."},
					},
				},
			},
		},
	}
}

func (d *paymentMethodsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *paymentMethodsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	methods, err := d.client.ListPaymentMethods(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing payment methods", err.Error())
		return
	}

	items := make([]paymentMethodItem, 0, len(methods))
	for _, m := range methods {
		items = append(items, paymentMethodItem{
			ID:   common.StringOrNull(m.ID),
			Name: common.StringOrNull(m.Name),
		})
	}

	data := paymentMethodsModel{
		ID:      types.StringValue("payment_methods"),
		Methods: items,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
