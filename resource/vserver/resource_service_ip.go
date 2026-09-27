// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	_ resource.Resource                = (*serviceIPResource)(nil)
	_ resource.ResourceWithConfigure   = (*serviceIPResource)(nil)
	_ resource.ResourceWithImportState = (*serviceIPResource)(nil)
)

// ServiceIPResource returns the hitechcloud_service_ip resource.
func ServiceIPResource() resource.Resource {
	return &serviceIPResource{}
}

type serviceIPResource struct {
	client *client.Client
}

type serviceIPModel struct {
	ID        types.String `tfsdk:"id"`
	ServiceID types.String `tfsdk:"service_id"`
	IPID      types.String `tfsdk:"ip_id"`
	VLAN      types.String `tfsdk:"vlan"`
	Domain    types.String `tfsdk:"domain"`
	IPAddress types.String `tfsdk:"ip_address"`
}

func (r *serviceIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_ip"
}

func (r *serviceIPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an IP address on a HiTechCloud bare metal / colocation " +
			"service (`POST /api/service/{service_id}/ips`). One resource manages exactly " +
			"one IP address.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `<service_id>/<ip_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the bare metal service. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ip_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "IP record identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vlan": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "VLAN to allocate the IP from. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"domain": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Reverse DNS domain assigned to the IP (editable in place).",
			},
			"ip_address": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The allocated IP address.",
			},
		},
	}
}

func (r *serviceIPResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *serviceIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serviceIPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ipID, err := r.client.CreateServiceIP(ctx, plan.ServiceID.ValueString(), plan.VLAN.ValueString(), plan.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating service IP", err.Error())
		return
	}

	plan.IPID = types.StringValue(ipID)
	plan.ID = types.StringValue(client.FormatID(plan.ServiceID.ValueString(), ipID))
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *serviceIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serviceIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, ipID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}
	state.ServiceID = types.StringValue(serviceID)
	state.IPID = types.StringValue(ipID)
	if !r.refresh(ctx, &state, &resp.Diagnostics) {
		resp.State.RemoveResource(ctx)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *serviceIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan serviceIPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, ipID, ok := r.parseID(ctx, &plan, &resp.Diagnostics)
	if !ok {
		return
	}

	if !plan.Domain.IsNull() && !plan.Domain.IsUnknown() {
		if err := r.client.UpdateServiceIP(ctx, serviceID, ipID, plan.Domain.ValueString()); err != nil {
			resp.Diagnostics.AddError("Error updating service IP", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(client.FormatID(serviceID, ipID))
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *serviceIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serviceIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, ipID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	if err := r.client.DeleteServiceIP(ctx, serviceID, ipID); client.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error deleting service IP", err.Error())
	}
}

func (r *serviceIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := client.SplitID(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ip_id"), parts[1])...)
}

func (r *serviceIPResource) parseID(_ context.Context, state *serviceIPModel, diags *diag.Diagnostics) (string, string, bool) {
	if state.ServiceID.ValueString() != "" && state.IPID.ValueString() != "" {
		return state.ServiceID.ValueString(), state.IPID.ValueString(), true
	}
	parts, err := client.SplitID(state.ID.ValueString(), 2)
	if err != nil {
		diags.AddError("Unexpected Resource ID", err.Error())
		return "", "", false
	}
	return parts[0], parts[1], true
}

func (r *serviceIPResource) refresh(ctx context.Context, model *serviceIPModel, diags *diag.Diagnostics) bool {
	ip, err := r.client.GetServiceIP(ctx, model.ServiceID.ValueString(), model.IPID.ValueString())
	if client.IsNotFound(err) {
		return false
	}
	if err != nil {
		diags.AddError("Error reading service IP", err.Error())
		return true
	}
	model.VLAN = common.StringOrNull(ip.VLAN)
	model.Domain = common.StringOrNull(ip.Domain)
	model.IPAddress = common.StringOrNull(ip.IP)
	return true
}
