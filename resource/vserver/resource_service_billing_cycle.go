// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vserver

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
	_ resource.Resource                = (*serviceBillingCycleResource)(nil)
	_ resource.ResourceWithConfigure   = (*serviceBillingCycleResource)(nil)
	_ resource.ResourceWithImportState = (*serviceBillingCycleResource)(nil)
)

// ServiceBillingCycleResource returns the hitechcloud_service_billing_cycle resource.
func ServiceBillingCycleResource() resource.Resource {
	return &serviceBillingCycleResource{}
}

type serviceBillingCycleResource struct {
	client *client.Client
}

type serviceBillingCycleModel struct {
	ID        types.String `tfsdk:"id"`
	ServiceID types.String `tfsdk:"service_id"`
	Cycle     types.String `tfsdk:"cycle"`
}

func (r *serviceBillingCycleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_billing_cycle"
}

func (r *serviceBillingCycleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the billing cycle of a HiTechCloud service " +
			"(`POST /api/service/{id}/cycle`). Destroying the resource leaves the " +
			"cycle unchanged (the API has no reset).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The service ID (same as `service_id`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the service. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cycle": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Billing cycle, for example `monthly`, `quarterly`, `annually`.",
			},
		},
	}
}

func (r *serviceBillingCycleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *serviceBillingCycleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serviceBillingCycleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.ChangeServiceBillingCycle(ctx, plan.ServiceID.ValueString(), plan.Cycle.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error setting service billing cycle", err.Error())
		return
	}
	plan.ID = types.StringValue(plan.ServiceID.ValueString())
	r.refresh(ctx, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *serviceBillingCycleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serviceBillingCycleModel
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

func (r *serviceBillingCycleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan serviceBillingCycleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.ChangeServiceBillingCycle(ctx, plan.ServiceID.ValueString(), plan.Cycle.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error updating service billing cycle", err.Error())
		return
	}
	r.refresh(ctx, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *serviceBillingCycleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serviceBillingCycleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.AddWarning(
		"Billing cycle unchanged",
		"The HiTechCloud API has no billing cycle reset; the attribute was removed from Terraform state only.",
	)
}

func (r *serviceBillingCycleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), req.ID)...)
}

// refresh reloads the billing cycle and reports whether the service exists.
func (r *serviceBillingCycleResource) refresh(ctx context.Context, model *serviceBillingCycleModel) bool {
	raw, err := r.client.GetServiceBillingCycle(ctx, model.ServiceID.ValueString())
	if err != nil {
		return !client.IsNotFound(err)
	}
	if v := client.FirstString(raw, "cycle", "billing_cycle", "billingcycle"); v != "" {
		model.Cycle = common.StringOrNull(v)
	}
	return true
}
