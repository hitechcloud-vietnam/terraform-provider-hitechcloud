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
	_ datasource.DataSource              = (*productsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*productsDataSource)(nil)
)

// ProductsDataSource returns the hitechcloud_products data source.
func ProductsDataSource() datasource.DataSource {
	return &productsDataSource{}
}

type productsDataSource struct {
	client *client.Client
}

type productsModel struct {
	ID         types.String  `tfsdk:"id"`
	CategoryID types.String  `tfsdk:"category_id"`
	Products   []productItem `tfsdk:"products"`
}

type productItem struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Cycle       types.String `tfsdk:"cycle"`
	Price       types.String `tfsdk:"price"`
}

func (d *productsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_products"
}

func (d *productsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the purchasable products of a HiTechCloud product " +
			"category (`GET /api/category/{category_id}/product`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (same as `category_id`).",
			},
			"category_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Product category identifier.",
			},
			"products": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Products in the category.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true, MarkdownDescription: "Product identifier."},
						"name":        schema.StringAttribute{Computed: true, MarkdownDescription: "Product name."},
						"description": schema.StringAttribute{Computed: true, MarkdownDescription: "Product description."},
						"cycle":       schema.StringAttribute{Computed: true, MarkdownDescription: "Billing cycle."},
						"price":       schema.StringAttribute{Computed: true, MarkdownDescription: "Price."},
					},
				},
			},
		},
	}
}

func (d *productsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *productsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data productsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	categoryID := data.CategoryID.ValueString()
	products, err := d.client.ListProducts(ctx, categoryID)
	if err != nil {
		resp.Diagnostics.AddError("Error listing products", err.Error())
		return
	}

	items := make([]productItem, 0, len(products))
	for _, p := range products {
		items = append(items, productItem{
			ID:          common.StringOrNull(p.ID),
			Name:        common.StringOrNull(p.Name),
			Description: common.StringOrNull(p.Description),
			Cycle:       common.StringOrNull(p.Cycle),
			Price:       common.StringOrNull(p.Price),
		})
	}

	data.ID = types.StringValue(categoryID)
	data.Products = items
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
