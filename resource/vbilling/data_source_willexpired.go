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
	_ datasource.DataSource              = (*willExpiredDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*willExpiredDataSource)(nil)
)

// WillExpiredDataSource returns the hitechcloud_willexpired data source.
func WillExpiredDataSource() datasource.DataSource {
	return &willExpiredDataSource{}
}

type willExpiredDataSource struct {
	client *client.Client
}

type willExpiredModel struct {
	ID       types.String      `tfsdk:"id"`
	ItemType types.String      `tfsdk:"item_type"`
	Status   types.String      `tfsdk:"status"`
	Items    []willExpiredItem `tfsdk:"items"`
	Summary  types.Map         `tfsdk:"summary"`
	Invoices []willExpiredInv  `tfsdk:"upcoming_invoices"`
}

type willExpiredItem struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Type      types.String `tfsdk:"type"`
	ExpiresAt types.String `tfsdk:"expires_at"`
	Status    types.String `tfsdk:"status"`
	AutoRenew types.Bool   `tfsdk:"autorenew"`
}

type willExpiredInv struct {
	ID     types.String `tfsdk:"id"`
	ItemID types.String `tfsdk:"item_id"`
	Amount types.String `tfsdk:"amount"`
	DueAt  types.String `tfsdk:"due_at"`
}

func (d *willExpiredDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_willexpired"
}

func (d *willExpiredDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists services and domains approaching their expiry date " +
			"(`GET /api/willexpired`, `/willexpired/summary`, `/willexpired/invoices`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Static identifier for this query.",
			},
			"item_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter by item type: `service` or `domain`. Empty lists both.",
			},
			"status":  schema.StringAttribute{Optional: true, MarkdownDescription: "Filter by expiry status."},
			"summary": schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Aggregate counters of the expiring items."},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Items approaching expiry.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.StringAttribute{Computed: true, MarkdownDescription: "Item ID (service ID or domain name)."},
						"name":       schema.StringAttribute{Computed: true, MarkdownDescription: "Item display name."},
						"type":       schema.StringAttribute{Computed: true, MarkdownDescription: "Item type (service/domain)."},
						"expires_at": schema.StringAttribute{Computed: true, MarkdownDescription: "Expiry timestamp."},
						"status":     schema.StringAttribute{Computed: true, MarkdownDescription: "Expiry status."},
						"autorenew":  schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether auto-renew is enabled."},
					},
				},
			},
			"upcoming_invoices": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Invoices generated for expiring items.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":      schema.StringAttribute{Computed: true, MarkdownDescription: "Invoice ID."},
						"item_id": schema.StringAttribute{Computed: true, MarkdownDescription: "Related item ID."},
						"amount":  schema.StringAttribute{Computed: true, MarkdownDescription: "Invoice amount."},
						"due_at":  schema.StringAttribute{Computed: true, MarkdownDescription: "Invoice due date."},
					},
				},
			},
		},
	}
}

func (d *willExpiredDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *willExpiredDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data willExpiredModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]string{}
	if v := data.ItemType.ValueString(); v != "" {
		params["type"] = v
	}
	if v := data.Status.ValueString(); v != "" {
		params["status"] = v
	}
	data.ID = types.StringValue("willexpired")

	items, err := d.client.ListWillExpired(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Error reading expiring items", err.Error())
		return
	}
	data.Items = make([]willExpiredItem, 0, len(items))
	for _, item := range items {
		m, ok := client.AsMap(item)
		if !ok {
			continue
		}
		data.Items = append(data.Items, willExpiredItem{
			ID:        types.StringValue(client.FirstString(m, "id", "item_id")),
			Name:      common.StringOrNull(client.FirstString(m, "name", "domain", "title")),
			Type:      common.StringOrNull(client.FirstString(m, "type", "item_type")),
			ExpiresAt: common.StringOrNull(client.FirstString(m, "expires_at", "expire", "expiry", "due")),
			Status:    common.StringOrNull(client.FirstString(m, "status")),
			AutoRenew: types.BoolValue(client.FirstBool(m, "autorenew", "auto_renew")),
		})
	}

	if summary, err := d.client.GetWillExpiredSummary(ctx, params); err == nil {
		data.Summary = common.MapStrings(client.StringMap(summary))
	}

	if invoices, err := d.client.ListWillExpiredInvoices(ctx, params); err == nil {
		data.Invoices = make([]willExpiredInv, 0, len(invoices))
		for _, item := range invoices {
			m, ok := client.AsMap(item)
			if !ok {
				continue
			}
			data.Invoices = append(data.Invoices, willExpiredInv{
				ID:     common.StringOrNull(client.FirstString(m, "id", "invoice_id")),
				ItemID: common.StringOrNull(client.FirstString(m, "item_id", "service_id", "domain")),
				Amount: common.StringOrNull(client.FirstString(m, "amount", "total", "sum")),
				DueAt:  common.StringOrNull(client.FirstString(m, "due_at", "due", "date")),
			})
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
