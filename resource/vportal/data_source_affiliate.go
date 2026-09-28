// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vportal

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
	_ datasource.DataSource              = (*affiliateDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*affiliateDataSource)(nil)
)

// AffiliateDataSource returns the hitechcloud_affiliate data source.
func AffiliateDataSource() datasource.DataSource {
	return &affiliateDataSource{}
}

type affiliateDataSource struct {
	client *client.Client
}

type affiliateModel struct {
	ID              types.String   `tfsdk:"id"`
	Summary         types.Map      `tfsdk:"summary"`
	Campaigns       []affiliateRow `tfsdk:"campaigns"`
	Commissions     []affiliateRow `tfsdk:"commissions"`
	Payouts         []affiliateRow `tfsdk:"payouts"`
	Vouchers        []affiliateRow `tfsdk:"vouchers"`
	CommissionPlans []affiliateRow `tfsdk:"commission_plans"`
}

type affiliateRow struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Status types.String `tfsdk:"status"`
	Amount types.String `tfsdk:"amount"`
	Date   types.String `tfsdk:"date"`
	Extra  types.Map    `tfsdk:"extra"`
}

func (d *affiliateDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_affiliate"
}

func (d *affiliateDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads the caller's affiliate program data: summary, campaigns, " +
			"commissions, payouts, vouchers and commission plans " +
			"(`GET /api/affiliates/summary`, `/affiliates/campaigns`, `/affiliates/commissions`, " +
			"`/affiliates/payouts`, `/affiliates/vouchers`, `/affiliates/commissionplans`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Static identifier for this query (`affiliate`).",
			},
			"summary": schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Affiliate account summary counters and balances."},
			"campaigns": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Affiliate campaigns.",
				NestedObject:        affiliateRowSchema(),
			},
			"commissions": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Earned commissions.",
				NestedObject:        affiliateRowSchema(),
			},
			"payouts": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Payout history.",
				NestedObject:        affiliateRowSchema(),
			},
			"vouchers": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Affiliate vouchers.",
				NestedObject:        affiliateRowSchema(),
			},
			"commission_plans": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Available commission plans.",
				NestedObject:        affiliateRowSchema(),
			},
		},
	}
}

// affiliateRowSchema is the shared nested schema for the affiliate listings.
func affiliateRowSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true, MarkdownDescription: "Row identifier."},
			"name":   schema.StringAttribute{Computed: true, MarkdownDescription: "Display name of the row."},
			"status": schema.StringAttribute{Computed: true, MarkdownDescription: "Status of the row."},
			"amount": schema.StringAttribute{Computed: true, MarkdownDescription: "Monetary amount, when present."},
			"date":   schema.StringAttribute{Computed: true, MarkdownDescription: "Date of the row, when present."},
			"extra":  schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Remaining fields of the row."},
		},
	}
}

func (d *affiliateDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// affiliateRowFrom converts one API item into the shared row model.
func affiliateRowFrom(item any) affiliateRow {
	m, ok := client.AsMap(item)
	if !ok {
		m = map[string]any{}
	}
	return affiliateRow{
		ID:     common.StringOrNull(client.FirstString(m, "id", "uid", "hash")),
		Name:   common.StringOrNull(client.FirstString(m, "name", "title", "campaign", "code")),
		Status: common.StringOrNull(client.FirstString(m, "status", "state")),
		Amount: common.StringOrNull(client.FirstString(m, "amount", "sum", "total", "value")),
		Date:   common.StringOrNull(client.FirstString(m, "date", "created_at", "created", "time")),
		Extra:  common.MapStrings(client.StringMap(m)),
	}
}

func (d *affiliateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data affiliateModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID = types.StringValue("affiliate")

	if summary, err := d.client.GetAffiliateSummary(ctx); err == nil {
		data.Summary = common.MapStrings(client.StringMap(summary))
	}
	if list, err := d.client.ListAffiliateCampaigns(ctx); err == nil {
		data.Campaigns = make([]affiliateRow, 0, len(list))
		for _, item := range list {
			data.Campaigns = append(data.Campaigns, affiliateRowFrom(item))
		}
	}
	if list, err := d.client.ListAffiliateCommissions(ctx); err == nil {
		data.Commissions = make([]affiliateRow, 0, len(list))
		for _, item := range list {
			data.Commissions = append(data.Commissions, affiliateRowFrom(item))
		}
	}
	if list, err := d.client.ListAffiliatePayouts(ctx); err == nil {
		data.Payouts = make([]affiliateRow, 0, len(list))
		for _, item := range list {
			data.Payouts = append(data.Payouts, affiliateRowFrom(item))
		}
	}
	if list, err := d.client.ListAffiliateVouchers(ctx); err == nil {
		data.Vouchers = make([]affiliateRow, 0, len(list))
		for _, item := range list {
			data.Vouchers = append(data.Vouchers, affiliateRowFrom(item))
		}
	}
	if list, err := d.client.ListAffiliateCommissionPlans(ctx); err == nil {
		data.CommissionPlans = make([]affiliateRow, 0, len(list))
		for _, item := range list {
			data.CommissionPlans = append(data.CommissionPlans, affiliateRowFrom(item))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
