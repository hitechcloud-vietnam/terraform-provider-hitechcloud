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
	_ datasource.DataSource              = (*invoicesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*invoicesDataSource)(nil)
)

// InvoicesDataSource returns the hitechcloud_invoices data source.
func InvoicesDataSource() datasource.DataSource {
	return &invoicesDataSource{}
}

type invoicesDataSource struct {
	client *client.Client
}

type invoicesModel struct {
	ID       types.String  `tfsdk:"id"`
	Invoices []invoiceItem `tfsdk:"invoices"`
}

type invoiceItem struct {
	ID       types.String `tfsdk:"id"`
	Number   types.String `tfsdk:"number"`
	Status   types.String `tfsdk:"status"`
	Currency types.String `tfsdk:"currency"`
	Total    types.String `tfsdk:"total"`
	Date     types.String `tfsdk:"date"`
	DueDate  types.String `tfsdk:"due_date"`
}

func (d *invoicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_invoices"
}

func (d *invoicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the invoices of the HiTechCloud account (`GET /api/invoice`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (`invoices`).",
			},
			"invoices": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Invoices of the account.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":       schema.StringAttribute{Computed: true, MarkdownDescription: "Invoice identifier."},
						"number":   schema.StringAttribute{Computed: true, MarkdownDescription: "Invoice number."},
						"status":   schema.StringAttribute{Computed: true, MarkdownDescription: "Invoice status."},
						"currency": schema.StringAttribute{Computed: true, MarkdownDescription: "Currency."},
						"total":    schema.StringAttribute{Computed: true, MarkdownDescription: "Total amount due."},
						"date":     schema.StringAttribute{Computed: true, MarkdownDescription: "Issue date."},
						"due_date": schema.StringAttribute{Computed: true, MarkdownDescription: "Due date."},
					},
				},
			},
		},
	}
}

func (d *invoicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *invoicesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	invoices, err := d.client.ListInvoices(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing invoices", err.Error())
		return
	}

	items := make([]invoiceItem, 0, len(invoices))
	for _, inv := range invoices {
		items = append(items, invoiceItem{
			ID:       common.StringOrNull(inv.ID),
			Number:   common.StringOrNull(inv.Number),
			Status:   common.StringOrNull(inv.Status),
			Currency: common.StringOrNull(inv.Currency),
			Total:    common.StringOrNull(inv.Total),
			Date:     common.StringOrNull(inv.Date),
			DueDate:  common.StringOrNull(inv.DueDate),
		})
	}

	data := invoicesModel{
		ID:       types.StringValue("invoices"),
		Invoices: items,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
