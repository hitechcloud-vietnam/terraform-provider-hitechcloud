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
	_ datasource.DataSource              = (*affiliateAdvDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*affiliateAdvDataSource)(nil)
)

// AffiliateAdvDataSource returns the hitechcloud_affiliate_adv data source.
func AffiliateAdvDataSource() datasource.DataSource {
	return &affiliateAdvDataSource{}
}

type affiliateAdvDataSource struct {
	client *client.Client
}

type affiliateAdvModel struct {
	ID          types.String   `tfsdk:"id"`
	ClientID    types.String   `tfsdk:"client_id"`
	Info        types.Map      `tfsdk:"info"`
	Stats       types.Map      `tfsdk:"stats"`
	Referrals   []affiliateRow `tfsdk:"referrals"`
	Vouchers    []affiliateRow `tfsdk:"vouchers"`
	Commissions []affiliateRow `tfsdk:"commissions"`
}

func (d *affiliateAdvDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_affiliate_adv"
}

func (d *affiliateAdvDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads the advanced affiliate data of one client: profile, " +
			"statistics, referrals, vouchers and commissions " +
			"(`GET /api/affiliates_adv/{client_id}/info`, `/stats`, `/referrals`, " +
			"`/vouchers`, `/commissions`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The client ID (same as `client_id`).",
			},
			"client_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Affiliate client ID.",
			},
			"info":  schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Affiliate profile fields."},
			"stats": schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Aggregate affiliate statistics."},
			"referrals": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Referred customers.",
				NestedObject:        affiliateRowSchema(),
			},
			"vouchers": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Vouchers created under the affiliate.",
				NestedObject:        affiliateRowSchema(),
			},
			"commissions": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Commissions earned by the affiliate.",
				NestedObject:        affiliateRowSchema(),
			},
		},
	}
}

func (d *affiliateAdvDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *affiliateAdvDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data affiliateAdvModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientID := data.ClientID.ValueString()
	data.ID = types.StringValue(clientID)

	if info, err := d.client.GetAffiliateAdvInfo(ctx, clientID); err == nil {
		data.Info = common.MapStrings(client.StringMap(info))
	}
	if stats, err := d.client.GetAffiliateAdvStats(ctx, clientID); err == nil {
		data.Stats = common.MapStrings(client.StringMap(stats))
	}
	if list, err := d.client.ListAffiliateAdvReferrals(ctx, clientID); err == nil {
		data.Referrals = make([]affiliateRow, 0, len(list))
		for _, item := range list {
			data.Referrals = append(data.Referrals, affiliateRowFrom(item))
		}
	}
	if list, err := d.client.ListAffiliateAdvVouchers(ctx, clientID); err == nil {
		data.Vouchers = make([]affiliateRow, 0, len(list))
		for _, item := range list {
			data.Vouchers = append(data.Vouchers, affiliateRowFrom(item))
		}
	}
	if list, err := d.client.ListAffiliateAdvCommissions(ctx, clientID, map[string]string{}); err == nil {
		data.Commissions = make([]affiliateRow, 0, len(list))
		for _, item := range list {
			data.Commissions = append(data.Commissions, affiliateRowFrom(item))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
