// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vserver

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
	_ datasource.DataSource              = (*serverStockDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*serverStockDataSource)(nil)
)

// ServerStockDataSource returns the hitechcloud_server_stock data source.
func ServerStockDataSource() datasource.DataSource {
	return &serverStockDataSource{}
}

type serverStockDataSource struct {
	client *client.Client
}

type serverStockModel struct {
	ID    types.String     `tfsdk:"id"`
	Stock []serverStockRow `tfsdk:"stock"`
}

type serverStockRow struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Location types.String `tfsdk:"location"`
	Quantity types.Int64  `tfsdk:"quantity"`
	Extra    types.Map    `tfsdk:"extra"`
}

func (d *serverStockDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_stock"
}

func (d *serverStockDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists bare-metal server stock availability " +
			"(`GET /api/serverstock`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Static identifier for this query (`server_stock`).",
			},
			"stock": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Server stock entries.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":       schema.StringAttribute{Computed: true, MarkdownDescription: "Stock entry identifier."},
						"name":     schema.StringAttribute{Computed: true, MarkdownDescription: "Server model / product name."},
						"location": schema.StringAttribute{Computed: true, MarkdownDescription: "Datacenter location."},
						"quantity": schema.Int64Attribute{Computed: true, MarkdownDescription: "Units available."},
						"extra":    schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Remaining fields of the entry."},
					},
				},
			},
		},
	}
}

func (d *serverStockDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *serverStockDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data serverStockModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	stock, err := d.client.ListServerStock(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading server stock", err.Error())
		return
	}
	data.ID = types.StringValue("server_stock")
	data.Stock = make([]serverStockRow, 0, len(stock))
	for _, item := range stock {
		m, ok := client.AsMap(item)
		if !ok {
			continue
		}
		data.Stock = append(data.Stock, serverStockRow{
			ID:       common.StringOrNull(client.FirstString(m, "id", "sku", "code")),
			Name:     common.StringOrNull(client.FirstString(m, "name", "model", "product", "title")),
			Location: common.StringOrNull(client.FirstString(m, "location", "datacenter", "dc", "region")),
			Quantity: types.Int64Value(client.FirstInt64(m, "quantity", "qty", "count", "available")),
			Extra:    common.MapStrings(client.StringMap(m)),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
