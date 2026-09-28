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
	_ datasource.DataSource              = (*partnerDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*partnerDataSource)(nil)
)

// PartnerDataSource returns the hitechcloud_partner data source.
func PartnerDataSource() datasource.DataSource {
	return &partnerDataSource{}
}

type partnerDataSource struct {
	client *client.Client
}

type partnerModel struct {
	ID       types.String     `tfsdk:"id"`
	Profile  types.Map        `tfsdk:"profile"`
	Tiers    []partnerTier    `tfsdk:"tiers"`
	RateCard types.Map        `tfsdk:"rate_card"`
	Wallet   types.Map        `tfsdk:"wallet"`
	Leads    []partnerLeadRow `tfsdk:"recent_leads"`
}

type partnerTier struct {
	Name        types.String `tfsdk:"name"`
	Level       types.Int64  `tfsdk:"level"`
	Commission  types.String `tfsdk:"commission"`
	Requirement types.String `tfsdk:"requirement"`
}

type partnerLeadRow struct {
	ID        types.String `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	Company   types.String `tfsdk:"company"`
	Status    types.String `tfsdk:"status"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func (d *partnerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_partner"
}

func (d *partnerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads the caller's HiTechCloud partner program data: profile, " +
			"tiers, rate card, wallet and recent leads " +
			"(`GET /api/partner`, `/partner/tiers`, `/partner/ratecard`, `/partner/wallet`, `/partner/leads`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Partner identifier (or email) reported by the API.",
			},
			"profile":   schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Partner profile fields."},
			"rate_card": schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Commission rate card values."},
			"wallet":    schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Partner wallet balances and payout settings."},
			"tiers": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Partner program tiers and their requirements.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":        schema.StringAttribute{Computed: true, MarkdownDescription: "Tier name."},
						"level":       schema.Int64Attribute{Computed: true, MarkdownDescription: "Tier level."},
						"commission":  schema.StringAttribute{Computed: true, MarkdownDescription: "Commission rate of the tier."},
						"requirement": schema.StringAttribute{Computed: true, MarkdownDescription: "Requirement description."},
					},
				},
			},
			"recent_leads": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Most recent partner leads registered by the caller.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.StringAttribute{Computed: true, MarkdownDescription: "Lead identifier."},
						"email":      schema.StringAttribute{Computed: true, MarkdownDescription: "Lead email."},
						"company":    schema.StringAttribute{Computed: true, MarkdownDescription: "Lead company."},
						"status":     schema.StringAttribute{Computed: true, MarkdownDescription: "Lead status."},
						"created_at": schema.StringAttribute{Computed: true, MarkdownDescription: "Registration timestamp."},
					},
				},
			},
		},
	}
}

func (d *partnerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *partnerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data partnerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if profile, err := d.client.GetPartnerProfile(ctx); err == nil {
		data.Profile = common.MapStrings(client.StringMap(profile))
		data.ID = types.StringValue(client.FirstString(profile, "id", "partner_id", "email"))
	}
	if data.ID.IsNull() || data.ID.ValueString() == "" {
		data.ID = types.StringValue("partner")
	}

	if tiers, err := d.client.ListPartnerTiers(ctx); err == nil {
		data.Tiers = make([]partnerTier, 0, len(tiers))
		for _, item := range tiers {
			m, ok := client.AsMap(item)
			if !ok {
				continue
			}
			data.Tiers = append(data.Tiers, partnerTier{
				Name:        types.StringValue(client.FirstString(m, "name", "tier", "title")),
				Level:       types.Int64Value(client.FirstInt64(m, "level", "rank")),
				Commission:  common.StringOrNull(client.FirstString(m, "commission", "rate", "percent")),
				Requirement: common.StringOrNull(client.FirstString(m, "requirement", "require", "description")),
			})
		}
	}

	if rates, err := d.client.GetPartnerRateCard(ctx); err == nil {
		data.RateCard = common.MapStrings(client.StringMap(rates))
	}
	if wallet, err := d.client.GetPartnerWallet(ctx, map[string]string{}); err == nil {
		data.Wallet = common.MapStrings(client.StringMap(wallet))
	}

	if leads, err := d.client.ListPartnerLeads(ctx, map[string]string{"limit": "10"}); err == nil {
		data.Leads = make([]partnerLeadRow, 0, len(leads))
		for _, item := range leads {
			m, ok := client.AsMap(item)
			if !ok {
				continue
			}
			data.Leads = append(data.Leads, partnerLeadRow{
				ID:        common.StringOrNull(client.FirstString(m, "id", "lead_id")),
				Email:     common.StringOrNull(client.FirstString(m, "email")),
				Company:   common.StringOrNull(client.FirstString(m, "company")),
				Status:    common.StringOrNull(client.FirstString(m, "status")),
				CreatedAt: common.StringOrNull(client.FirstString(m, "created_at", "created")),
			})
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
