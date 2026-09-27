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
	_ datasource.DataSource              = (*categoriesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*categoriesDataSource)(nil)
)

// CategoriesDataSource returns the hitechcloud_categories data source.
func CategoriesDataSource() datasource.DataSource {
	return &categoriesDataSource{}
}

type categoriesDataSource struct {
	client *client.Client
}

type categoriesModel struct {
	ID         types.String   `tfsdk:"id"`
	Categories []categoryItem `tfsdk:"categories"`
}

type categoryItem struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (d *categoriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_categories"
}

func (d *categoriesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the product categories of the HiTechCloud catalog " +
			"(`GET /api/category`). Useful to discover the `category_id` used by " +
			"the `hitechcloud_products` data source.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (`categories`).",
			},
			"categories": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Product categories.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true, MarkdownDescription: "Category identifier."},
						"name":        schema.StringAttribute{Computed: true, MarkdownDescription: "Category name."},
						"description": schema.StringAttribute{Computed: true, MarkdownDescription: "Category description."},
					},
				},
			},
		},
	}
}

func (d *categoriesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *categoriesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	cats, err := d.client.ListProductCategories(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing product categories", err.Error())
		return
	}

	items := make([]categoryItem, 0, len(cats))
	for _, c := range cats {
		items = append(items, categoryItem{
			ID:          common.StringOrNull(c.ID),
			Name:        common.StringOrNull(c.Name),
			Description: common.StringOrNull(c.Description),
		})
	}

	data := categoriesModel{
		ID:         types.StringValue("categories"),
		Categories: items,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
