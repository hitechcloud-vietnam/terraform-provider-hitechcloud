// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vdomain

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
	_ resource.Resource                = (*emailForwardingResource)(nil)
	_ resource.ResourceWithConfigure   = (*emailForwardingResource)(nil)
	_ resource.ResourceWithImportState = (*emailForwardingResource)(nil)
)

// EmailForwardingResource returns the hitechcloud_email_forwarding resource.
func EmailForwardingResource() resource.Resource {
	return &emailForwardingResource{}
}

type emailForwardingResource struct {
	client *client.Client
}

type emailForwardingModel struct {
	ID       types.String `tfsdk:"id"`
	DomainID types.String `tfsdk:"domain_id"`
	From     types.String `tfsdk:"from"`
	To       types.String `tfsdk:"to"`
}

func (r *emailForwardingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_email_forwarding"
}

func (r *emailForwardingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an email forwarding rule of a registered domain " +
			"(`PUT /api/domain/{id}/emforwarding`). Destroying the resource clears the " +
			"forwarding rule.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `domain_id/from`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the registered domain. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"from": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Source email address to forward. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"to": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Destination email address receiving forwarded mail.",
			},
		},
	}
}

func (r *emailForwardingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *emailForwardingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan emailForwardingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateEmailForwarding(ctx,
		plan.DomainID.ValueString(),
		plan.From.ValueString(),
		plan.To.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("Error setting email forwarding", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.DomainID.ValueString(), plan.From.ValueString()))
	r.refresh(ctx, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *emailForwardingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state emailForwardingModel
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

func (r *emailForwardingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan emailForwardingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateEmailForwarding(ctx,
		plan.DomainID.ValueString(),
		plan.From.ValueString(),
		plan.To.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("Error updating email forwarding", err.Error())
		return
	}

	r.refresh(ctx, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *emailForwardingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state emailForwardingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// An empty destination clears the forwarding rule.
	if err := r.client.UpdateEmailForwarding(ctx,
		state.DomainID.ValueString(),
		state.From.ValueString(),
		"",
	); err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error clearing email forwarding", err.Error())
	}
}

func (r *emailForwardingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	parts, err := client.SplitID(id, 2)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf(
			"Expected import ID in the format domain_id/from: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("from"), parts[1])...)
}

// refresh reloads the forwarding rules and reports whether the rule exists.
func (r *emailForwardingResource) refresh(ctx context.Context, model *emailForwardingModel) bool {
	rules, err := r.client.GetEmailForwarding(ctx, model.DomainID.ValueString())
	if err != nil {
		return !client.IsNotFound(err)
	}
	// The response is either {from: to} pairs or a wrapped list of records.
	if to, ok := client.StringMap(rules)[model.From.ValueString()]; ok {
		model.To = common.StringOrNull(to)
		return true
	}
	if list := client.ExtractList(rules); len(list) > 0 {
		for _, item := range list {
			m, ok := client.AsMap(item)
			if !ok {
				continue
			}
			if client.FirstString(m, "from", "source") == model.From.ValueString() {
				model.To = common.StringOrNull(client.FirstString(m, "to", "destination", "target"))
				return true
			}
		}
	}
	return false
}
