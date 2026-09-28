// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vbilling

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
)

var (
	_ resource.Resource                = (*serviceAutoRenewResource)(nil)
	_ resource.ResourceWithConfigure   = (*serviceAutoRenewResource)(nil)
	_ resource.ResourceWithImportState = (*serviceAutoRenewResource)(nil)
)

// ServiceAutoRenewResource returns the hitechcloud_service_autorenew resource.
func ServiceAutoRenewResource() resource.Resource {
	return &serviceAutoRenewResource{}
}

type serviceAutoRenewResource struct {
	client *client.Client
}

type serviceAutoRenewModel struct {
	ID        types.String `tfsdk:"id"`
	ItemType  types.String `tfsdk:"item_type"`
	ItemID    types.String `tfsdk:"item_id"`
	AutoRenew types.Bool   `tfsdk:"autorenew"`
}

func (r *serviceAutoRenewResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_autorenew"
}

func (r *serviceAutoRenewResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the automatic-renewal flag of a service or domain " +
			"(`PUT /api/willexpired/{type}/{id}/autorenew`). Destroying the resource " +
			"turns auto-renew off.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `item_type/item_id`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"item_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Item type: `service` or `domain`. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"item_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Service ID (`hb_accounts.id`) or domain name (`hb_domains.id`). Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"autorenew": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether automatic renewal is enabled for the item.",
			},
		},
	}
}

func (r *serviceAutoRenewResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *serviceAutoRenewResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serviceAutoRenewModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.apply(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error setting auto-renew", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.ItemType.ValueString(), plan.ItemID.ValueString()))
	r.refresh(ctx, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *serviceAutoRenewResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serviceAutoRenewModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !r.refresh(ctx, &state) {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *serviceAutoRenewResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan serviceAutoRenewModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.apply(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error updating auto-renew", err.Error())
		return
	}

	r.refresh(ctx, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *serviceAutoRenewResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serviceAutoRenewModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.AutoRenew = types.BoolValue(false)
	if err := r.apply(ctx, &state); err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error disabling auto-renew", err.Error())
	}
}

func (r *serviceAutoRenewResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	parts, err := client.SplitID(id, 2)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf(
			"Expected import ID in the format item_type/item_id: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("item_type"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("item_id"), parts[1])...)
}

func (r *serviceAutoRenewResource) apply(ctx context.Context, model *serviceAutoRenewModel) error {
	value := "0"
	if model.AutoRenew.ValueBool() {
		value = "1"
	}
	return r.client.SetWillExpiredAutoRenew(ctx,
		model.ItemType.ValueString(),
		model.ItemID.ValueString(),
		value,
	)
}

// refresh reloads the auto-renew state and reports whether the item exists.
func (r *serviceAutoRenewResource) refresh(ctx context.Context, model *serviceAutoRenewModel) bool {
	raw, err := r.client.GetWillExpiredAutoRenew(ctx,
		model.ItemType.ValueString(),
		model.ItemID.ValueString(),
	)
	if err != nil {
		return !client.IsNotFound(err)
	}
	model.AutoRenew = types.BoolValue(client.FirstBool(raw, "autorenew", "auto_renew", "status"))
	return true
}
