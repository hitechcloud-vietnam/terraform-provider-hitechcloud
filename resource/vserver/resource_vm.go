// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vserver

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/client"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/common"
)

var (
	_ resource.Resource                = (*vmResource)(nil)
	_ resource.ResourceWithConfigure   = (*vmResource)(nil)
	_ resource.ResourceWithImportState = (*vmResource)(nil)
)

// VMResource returns the hitechcloud_vm resource.
func VMResource() resource.Resource {
	return &vmResource{}
}

type vmResource struct {
	client *client.Client
}

type vmModel struct {
	ID          types.String `tfsdk:"id"`
	ServiceID   types.String `tfsdk:"service_id"`
	VMID        types.String `tfsdk:"vm_id"`
	Label       types.String `tfsdk:"label"`
	Status      types.String `tfsdk:"status"`
	Hostname    types.String `tfsdk:"hostname"`
	Note        types.String `tfsdk:"note"`
	TemplateID  types.String `tfsdk:"template_id"`
	Password    types.String `tfsdk:"password"`
	Memory      types.Int64  `tfsdk:"memory"`
	CPU         types.Int64  `tfsdk:"cpu"`
	CPUShare    types.Int64  `tfsdk:"cpu_share"`
	Disk        types.Int64  `tfsdk:"disk"`
	Swap        types.Int64  `tfsdk:"swap"`
	LicenseKey  types.String `tfsdk:"license_key"`
	LicenseType types.String `tfsdk:"license_type"`
	IPs         types.Set    `tfsdk:"ips"`
}

func (r *vmResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm"
}

func (r *vmResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a HiTechCloud virtual machine (Cloud Instance) " +
			"(`POST /api/service/{service_id}/vms`). Memory, CPU and CPU share can be " +
			"resized in place; all other changes force replacement. Deletion waits " +
			"(up to 30 minutes) until the VM is actually gone.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `<service_id>/<vm_id>`.",
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
				Computed:            true,
				MarkdownDescription: "Virtual machine identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"label": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Display label of the VM. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current VM status.",
			},
			"hostname": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Hostname assigned to the VM.",
			},
			"note": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Note attached to the VM. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"template_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "OS template id. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Root/administrator password. Forces replacement; never stored in API responses.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"memory": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Memory in MB (resizable in place).",
			},
			"cpu": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Number of CPU cores (resizable in place).",
			},
			"cpu_share": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "CPU share in percent (resizable in place).",
			},
			"disk": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Disk size in GB. Forces replacement.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"swap": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Swap size in GB. Forces replacement.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"license_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "License key for the operating system. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"license_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "License type. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ips": schema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "IP addresses assigned to the VM.",
			},
		},
	}
}

func (r *vmResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *vmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vmModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := client.VMCreate{
		Label:       plan.Label.ValueString(),
		TemplateID:  plan.TemplateID.ValueString(),
		Password:    plan.Password.ValueString(),
		Memory:      plan.Memory.ValueInt64(),
		CPU:         plan.CPU.ValueInt64(),
		CPUShare:    plan.CPUShare.ValueInt64(),
		Disk:        plan.Disk.ValueInt64(),
		Swap:        plan.Swap.ValueInt64(),
		Note:        plan.Note.ValueString(),
		LicenseKey:  plan.LicenseKey.ValueString(),
		LicenseType: plan.LicenseType.ValueString(),
	}
	vmID, err := r.client.CreateVM(ctx, plan.ServiceID.ValueString(), in)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VM", err.Error())
		return
	}

	plan.VMID = types.StringValue(vmID)
	plan.ID = types.StringValue(client.FormatID(plan.ServiceID.ValueString(), vmID))
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *vmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vmModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, vmID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}
	state.ServiceID = types.StringValue(serviceID)
	state.VMID = types.StringValue(vmID)
	if !r.refresh(ctx, &state, &resp.Diagnostics) {
		resp.State.RemoveResource(ctx)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *vmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan vmModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, vmID, ok := r.parseID(ctx, &plan, &resp.Diagnostics)
	if !ok {
		return
	}

	if err := r.client.UpdateVM(ctx, serviceID, vmID,
		plan.Memory.ValueInt64(), plan.CPU.ValueInt64(), plan.CPUShare.ValueInt64()); err != nil {
		resp.Diagnostics.AddError("Error updating VM", err.Error())
		return
	}

	plan.ID = types.StringValue(client.FormatID(serviceID, vmID))
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *vmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vmModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, vmID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	if err := r.client.DeleteVM(ctx, serviceID, vmID); client.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error deleting VM", err.Error())
		return
	}
	if err := r.client.WaitForVMDeleted(ctx, serviceID, vmID, vmDeleteTimeout); err != nil {
		resp.Diagnostics.AddError("Error waiting for VM deletion", err.Error())
	}
}

// vmDeleteTimeout bounds the deletion polling.
const vmDeleteTimeout = 30 * time.Minute

func (r *vmResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := client.SplitID(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vm_id"), parts[1])...)
}

func (r *vmResource) parseID(_ context.Context, state *vmModel, diags *diag.Diagnostics) (string, string, bool) {
	if state.ServiceID.ValueString() != "" && state.VMID.ValueString() != "" {
		return state.ServiceID.ValueString(), state.VMID.ValueString(), true
	}
	parts, err := client.SplitID(state.ID.ValueString(), 2)
	if err != nil {
		diags.AddError("Unexpected Resource ID", err.Error())
		return "", "", false
	}
	return parts[0], parts[1], true
}

func (r *vmResource) refresh(ctx context.Context, model *vmModel, diags *diag.Diagnostics) bool {
	vm, err := r.client.GetVM(ctx, model.ServiceID.ValueString(), model.VMID.ValueString())
	if client.IsNotFound(err) {
		return false
	}
	if err != nil {
		diags.AddError("Error reading VM", err.Error())
		return true
	}
	model.Label = common.StringOrNull(vm.Label)
	model.Status = common.StringOrNull(vm.Status)
	model.Hostname = common.StringOrNull(vm.Hostname)
	if !model.Note.IsNull() {
		model.Note = common.StringOrNull(vm.Note)
	}
	model.TemplateID = common.StringOrNull(vm.TemplateID)
	model.Memory = common.Int64OrNull(vm.Memory)
	model.CPU = common.Int64OrNull(vm.CPU)
	model.CPUShare = common.Int64OrNull(vm.CPUShare)
	model.Disk = common.Int64OrNull(vm.Disk)
	model.Swap = common.Int64OrNull(vm.Swap)
	model.IPs = common.SetStrings(vm.IPs)
	return true
}
