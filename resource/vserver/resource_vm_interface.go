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
	_ resource.Resource                = (*vmInterfaceResource)(nil)
	_ resource.ResourceWithConfigure   = (*vmInterfaceResource)(nil)
	_ resource.ResourceWithImportState = (*vmInterfaceResource)(nil)
)

// VMInterfaceResource returns the hitechcloud_vm_interface resource.
func VMInterfaceResource() resource.Resource {
	return &vmInterfaceResource{}
}

type vmInterfaceResource struct {
	client *client.Client
}

type vmInterfaceModel struct {
	ID        types.String `tfsdk:"id"`
	ServiceID types.String `tfsdk:"service_id"`
	VMID      types.String `tfsdk:"vm_id"`
	Interface types.String `tfsdk:"interface_id"`
	Bridge    types.String `tfsdk:"bridge"`
	Firewall  types.Bool   `tfsdk:"firewall"`
	IPv4      []string     `tfsdk:"ipv4"`
	IPv6      []string     `tfsdk:"ipv6"`
}

func (r *vmInterfaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm_interface"
}

func (r *vmInterfaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a network interface of a HiTechCloud virtual machine " +
			"(`POST /api/service/{service_id}/vms/{vm_id}/interfaces`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `<service_id>/<vm_id>/<interface_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the cloud service. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vm_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Virtual machine identifier. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"interface_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Interface identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bridge": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Network bridge to attach to. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"firewall": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the firewall is enabled on the interface.",
			},
			"ipv4": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "IPv4 address ids to assign.",
			},
			"ipv6": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "IPv6 address ids to assign.",
			},
		},
	}
}

func (r *vmInterfaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *vmInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vmInterfaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := client.VMInterface{
		Bridge:   plan.Bridge.ValueString(),
		Firewall: plan.Firewall.ValueBool(),
		IPv4IDs:  plan.IPv4,
		IPv6IDs:  plan.IPv6,
	}
	iface, err := r.client.CreateVMInterface(ctx, plan.ServiceID.ValueString(), plan.VMID.ValueString(), in)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VM interface", err.Error())
		return
	}

	plan.Interface = types.StringValue(iface)
	plan.ID = types.StringValue(client.FormatID(plan.ServiceID.ValueString(), plan.VMID.ValueString(), iface))
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *vmInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vmInterfaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, vmID, iface, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}
	state.ServiceID = types.StringValue(serviceID)
	state.VMID = types.StringValue(vmID)
	state.Interface = types.StringValue(iface)
	if !r.refresh(ctx, &state, &resp.Diagnostics) {
		resp.State.RemoveResource(ctx)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *vmInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan vmInterfaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, vmID, iface, ok := r.parseID(ctx, &plan, &resp.Diagnostics)
	if !ok {
		return
	}

	in := client.VMInterface{
		Firewall: plan.Firewall.ValueBool(),
		IPv4IDs:  plan.IPv4,
		IPv6IDs:  plan.IPv6,
	}
	if err := r.client.UpdateVMInterface(ctx, serviceID, vmID, iface, in); err != nil {
		resp.Diagnostics.AddError("Error updating VM interface", err.Error())
		return
	}

	plan.ID = types.StringValue(client.FormatID(serviceID, vmID, iface))
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *vmInterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vmInterfaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, vmID, iface, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	if err := r.client.DeleteVMInterface(ctx, serviceID, vmID, iface); client.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error deleting VM interface", err.Error())
	}
}

func (r *vmInterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := client.SplitID(req.ID, 3)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vm_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("interface_id"), parts[2])...)
}

func (r *vmInterfaceResource) parseID(_ context.Context, state *vmInterfaceModel, diags *diag.Diagnostics) (string, string, string, bool) {
	if state.ServiceID.ValueString() != "" && state.VMID.ValueString() != "" && state.Interface.ValueString() != "" {
		return state.ServiceID.ValueString(), state.VMID.ValueString(), state.Interface.ValueString(), true
	}
	parts, err := client.SplitID(state.ID.ValueString(), 3)
	if err != nil {
		diags.AddError("Unexpected Resource ID", err.Error())
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

func (r *vmInterfaceResource) refresh(ctx context.Context, model *vmInterfaceModel, diags *diag.Diagnostics) bool {
	nic, err := r.client.GetVMInterface(ctx, model.ServiceID.ValueString(), model.VMID.ValueString(), model.Interface.ValueString())
	if client.IsNotFound(err) {
		return false
	}
	if err != nil {
		diags.AddError("Error reading VM interface", err.Error())
		return true
	}
	model.Bridge = common.StringOrNull(nic.Bridge)
	model.Firewall = types.BoolValue(nic.Firewall)
	if len(nic.IPv4IDs) > 0 {
		model.IPv4 = nic.IPv4IDs
	} else if model.IPv4 != nil {
		model.IPv4 = []string{}
	}
	if len(nic.IPv6IDs) > 0 {
		model.IPv6 = nic.IPv6IDs
	} else if model.IPv6 != nil {
		model.IPv6 = []string{}
	}
	return true
}
