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
	_ datasource.DataSource              = (*productConfigDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*productConfigDataSource)(nil)
)

// ProductConfigDataSource returns the hitechcloud_product_config data source.
func ProductConfigDataSource() datasource.DataSource {
	return &productConfigDataSource{}
}

type productConfigDataSource struct {
	client *client.Client
}

type productConfigModel struct {
	ID        types.String `tfsdk:"id"`
	ProductID types.String `tfsdk:"product_id"`
	Config    types.Map    `tfsdk:"config"`
}

func (d *productConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_product_config"
}

func (d *productConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads the order form / configuration options of a product " +
			"(`GET /api/order/{product_id}`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The product ID (same as `product_id`).",
			},
			"product_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Product identifier from `hitechcloud_products`.",
			},
			"config": schema.MapAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Order form fields and available options of the product.",
			},
		},
	}
}

func (d *productConfigDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *productConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data productConfigModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	productID := data.ProductID.ValueString()
	data.ID = types.StringValue(productID)

	cfg, err := d.client.GetProductConfig(ctx, productID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading product config", err.Error())
		return
	}
	data.Config = common.MapStrings(client.StringMap(cfg))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
