// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vportal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/client"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/common"
)

var (
	_ resource.Resource                = (*partnerPayoutResource)(nil)
	_ resource.ResourceWithConfigure   = (*partnerPayoutResource)(nil)
	_ resource.ResourceWithImportState = (*partnerPayoutResource)(nil)
)

// PartnerPayoutResource returns the hitechcloud_partner_payout resource.
func PartnerPayoutResource() resource.Resource {
	return &partnerPayoutResource{}
}

type partnerPayoutResource struct {
	client *client.Client
}

type partnerPayoutModel struct {
	ID        types.String `tfsdk:"id"`
	Amount    types.String `tfsdk:"amount"`
	Method    types.String `tfsdk:"method"`
	Note      types.String `tfsdk:"note"`
	Status    types.String `tfsdk:"status"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func (r *partnerPayoutResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_partner_payout"
}

func (r *partnerPayoutResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Requests a partner wallet payout " +
			"(`POST /api/partner/payouts`). Payout requests are immutable: changing " +
			"any attribute forces a replacement. Deleting the resource removes it from " +
			"Terraform state only (the API has no payout cancellation endpoint).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Payout request identifier returned by the API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"amount": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Amount to pay out. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"method": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Payout method (e.g. bank transfer). Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"note": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional note for the payout request. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current payout status reported by the API.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Request timestamp reported by the API.",
			},
		},
	}
}

func (r *partnerPayoutResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	cli, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got %T", req.ProviderData),
		)
		return
	}
	r.client = cli
}

func (r *partnerPayoutResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan partnerPayoutModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payout, err := r.client.RequestPartnerPayout(ctx,
		plan.Amount.ValueString(),
		plan.Method.ValueString(),
		plan.Note.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error requesting partner payout", err.Error())
		return
	}

	id := client.FirstString(payout, "id", "payout_id", "request_id")
	if id == "" {
		id = plan.Amount.ValueString()
	}
	plan.ID = types.StringValue(id)
	plan.Status = common.StringOrNull(client.FirstString(payout, "status"))
	plan.CreatedAt = common.StringOrNull(client.FirstString(payout, "created_at", "created"))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *partnerPayoutResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state partnerPayoutModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payouts, err := r.client.ListPartnerPayouts(ctx, map[string]string{})
	if err != nil {
		resp.Diagnostics.AddError("Error reading partner payouts", err.Error())
		return
	}
	m, ok := client.FindByID(payouts, state.ID.ValueString(), "id", "payout_id", "request_id")
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}
	state.Status = common.StringOrNull(client.FirstString(m, "status"))
	state.CreatedAt = common.StringOrNull(client.FirstString(m, "created_at", "created"))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Payout requests are immutable: Update is only reached via replacement.
func (r *partnerPayoutResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Partner payouts are immutable",
		"The API does not support editing payout requests; this resource only supports create and read.",
	)
}

func (r *partnerPayoutResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Payout request not cancelled",
		"The HiTechCloud API does not expose a payout cancellation endpoint; the request was removed from Terraform state only.",
	)
}

func (r *partnerPayoutResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
