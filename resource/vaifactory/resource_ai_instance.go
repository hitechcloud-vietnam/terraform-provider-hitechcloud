// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vaifactory

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/client"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/common"
)

const (
	aiInstanceCreateTimeout = 60 * time.Minute
	aiInstanceDeleteTimeout = 30 * time.Minute
)

var (
	_ resource.Resource                = (*aiInstanceResource)(nil)
	_ resource.ResourceWithConfigure   = (*aiInstanceResource)(nil)
	_ resource.ResourceWithImportState = (*aiInstanceResource)(nil)
)

// AIInstanceResource returns the hitechcloud_ai_instance resource.
func AIInstanceResource() resource.Resource {
	return &aiInstanceResource{}
}

type aiInstanceResource struct {
	client *client.Client
}

type aiInstanceModel struct {
	ID                  types.String `tfsdk:"id"`
	ServiceID           types.String `tfsdk:"service_id"`
	InstanceID          types.String `tfsdk:"instance_id"`
	Name                types.String `tfsdk:"name"`
	Cloud               types.String `tfsdk:"cloud"`
	Region              types.String `tfsdk:"region"`
	ShadeInstanceType   types.String `tfsdk:"shade_instance_type"`
	ShadeCloud          types.Bool   `tfsdk:"shade_cloud"`
	OS                  types.String `tfsdk:"os"`
	TemplateID          types.String `tfsdk:"template_id"`
	SSHKeyID            types.String `tfsdk:"ssh_key_id"`
	VolumeIDs           types.Set    `tfsdk:"volume_ids"`
	LaunchConfiguration types.String `tfsdk:"launch_configuration"`
	AutoDelete          types.Bool   `tfsdk:"auto_delete"`
	Alert               types.Bool   `tfsdk:"alert"`
	VolumeMount         types.String `tfsdk:"volume_mount"`
	Tags                types.Set    `tfsdk:"tags"`
	Envs                types.Map    `tfsdk:"envs"`
	Status              types.String `tfsdk:"status"`
	IP                  types.String `tfsdk:"ip"`
}

func (r *aiInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ai_instance"
}

func (r *aiInstanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a HiTechCloud AI Factory GPU instance " +
			"(`POST /api/service/{service_id}/instances`). Creation waits (up to 60 " +
			"minutes) until the instance leaves its transient provisioning state; " +
			"deletion waits (up to 30 minutes) until the instance is gone. Only `name`, " +
			"`auto_delete`, `alert` and `tags` can be updated in place.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `<service_id>/<instance_id>`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the AI Factory service. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"instance_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "GPU instance identifier.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Instance name.",
			},
			"cloud": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Cloud provider (e.g. `shade`). Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"region": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Region. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"shade_instance_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "GPU instance type (see `hitechcloud_ai_instance_types`). Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"shade_cloud": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether the instance runs on the Shade Cloud. Forces replacement.",
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"os": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Operating system / image. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"template_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Template id to launch from. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"ssh_key_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "SSH key id to inject. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"volume_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Volume ids to attach at creation. Forces replacement.",
				PlanModifiers:       []planmodifier.Set{setplanmodifier.RequiresReplace()},
			},
			"launch_configuration": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Launch configuration as a JSON string (passed to the API verbatim). Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"auto_delete": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the instance is deleted automatically when idle.",
			},
			"alert": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether alerts are enabled for the instance.",
			},
			"volume_mount": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Volume mount configuration as a JSON string (passed to the API verbatim). Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"tags": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Tags attached to the instance (updatable in place).",
			},
			"envs": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Environment variables passed to the instance. Forces replacement.",
				PlanModifiers:       []planmodifier.Map{mapplanmodifier.RequiresReplace()},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current instance status.",
			},
			"ip": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "IP address of the instance.",
			},
		},
	}
}

func (r *aiInstanceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *aiInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aiInstanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := client.AIInstanceCreate{
		Cloud:               plan.Cloud.ValueString(),
		Region:              plan.Region.ValueString(),
		ShadeInstanceType:   plan.ShadeInstanceType.ValueString(),
		Name:                plan.Name.ValueString(),
		OS:                  plan.OS.ValueString(),
		TemplateID:          plan.TemplateID.ValueString(),
		SSHKeyID:            plan.SSHKeyID.ValueString(),
		LaunchConfiguration: plan.LaunchConfiguration.ValueString(),
		VolumeMount:         plan.VolumeMount.ValueString(),
	}
	if !plan.ShadeCloud.IsNull() && !plan.ShadeCloud.IsUnknown() {
		in.ShadeCloud = common.BoolPtr(plan.ShadeCloud.ValueBool())
	}
	if !plan.AutoDelete.IsNull() && !plan.AutoDelete.IsUnknown() {
		in.AutoDelete = common.BoolPtr(plan.AutoDelete.ValueBool())
	}
	if !plan.Alert.IsNull() && !plan.Alert.IsUnknown() {
		in.Alert = common.BoolPtr(plan.Alert.ValueBool())
	}
	if !plan.VolumeIDs.IsNull() && !plan.VolumeIDs.IsUnknown() {
		ids, ds := common.StringsFromSet(ctx, plan.VolumeIDs)
		resp.Diagnostics.Append(ds...)
		in.VolumeIDs = ids
	}
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		tags, ds := common.StringsFromSet(ctx, plan.Tags)
		resp.Diagnostics.Append(ds...)
		in.Tags = tags
	}
	if !plan.Envs.IsNull() && !plan.Envs.IsUnknown() {
		envs, ds := common.StringsFromMap(ctx, plan.Envs)
		resp.Diagnostics.Append(ds...)
		in.Envs = envs
	}
	if resp.Diagnostics.HasError() {
		return
	}

	instanceID, err := r.client.CreateAIInstance(ctx, plan.ServiceID.ValueString(), in)
	if err != nil {
		resp.Diagnostics.AddError("Error creating AI instance", err.Error())
		return
	}
	plan.InstanceID = types.StringValue(instanceID)
	plan.ID = types.StringValue(client.FormatID(plan.ServiceID.ValueString(), instanceID))

	// GPU instances are provisioned asynchronously; wait for a stable state.
	if _, err := r.client.WaitForAIInstanceStable(ctx, plan.ServiceID.ValueString(), instanceID, aiInstanceCreateTimeout); err != nil {
		resp.Diagnostics.AddError("Error waiting for AI instance", err.Error())
		return
	}

	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *aiInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state aiInstanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, instanceID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}
	state.ServiceID = types.StringValue(serviceID)
	state.InstanceID = types.StringValue(instanceID)
	if !r.refresh(ctx, &state, &resp.Diagnostics) {
		resp.State.RemoveResource(ctx)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *aiInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan aiInstanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, instanceID, ok := r.parseID(ctx, &plan, &resp.Diagnostics)
	if !ok {
		return
	}

	var autoDelete, alert *bool
	if !plan.AutoDelete.IsNull() && !plan.AutoDelete.IsUnknown() {
		autoDelete = common.BoolPtr(plan.AutoDelete.ValueBool())
	}
	if !plan.Alert.IsNull() && !plan.Alert.IsUnknown() {
		alert = common.BoolPtr(plan.Alert.ValueBool())
	}
	var tags []string
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		tags, _ = common.StringsFromSet(ctx, plan.Tags)
	}

	if err := r.client.UpdateAIInstance(ctx, serviceID, instanceID,
		plan.Name.ValueString(), autoDelete, alert, tags); err != nil {
		resp.Diagnostics.AddError("Error updating AI instance", err.Error())
		return
	}

	plan.ID = types.StringValue(client.FormatID(serviceID, instanceID))
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *aiInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state aiInstanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, instanceID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	if err := r.client.DeleteAIInstance(ctx, serviceID, instanceID); client.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error deleting AI instance", err.Error())
		return
	}
	if err := r.client.WaitForAIInstanceDeleted(ctx, serviceID, instanceID, aiInstanceDeleteTimeout); err != nil {
		resp.Diagnostics.AddError("Error waiting for AI instance deletion", err.Error())
	}
}

func (r *aiInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := client.SplitID(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance_id"), parts[1])...)
}

func (r *aiInstanceResource) parseID(_ context.Context, state *aiInstanceModel, diags *diag.Diagnostics) (string, string, bool) {
	if state.ServiceID.ValueString() != "" && state.InstanceID.ValueString() != "" {
		return state.ServiceID.ValueString(), state.InstanceID.ValueString(), true
	}
	parts, err := client.SplitID(state.ID.ValueString(), 2)
	if err != nil {
		diags.AddError("Unexpected Resource ID", err.Error())
		return "", "", false
	}
	return parts[0], parts[1], true
}

// refresh reloads the instance; returns false when it no longer exists.
func (r *aiInstanceResource) refresh(ctx context.Context, model *aiInstanceModel, diags *diag.Diagnostics) bool {
	inst, err := r.client.GetAIInstance(ctx, model.ServiceID.ValueString(), model.InstanceID.ValueString())
	if client.IsNotFound(err) {
		return false
	}
	if err != nil {
		diags.AddError("Error reading AI instance", err.Error())
		return true
	}
	model.Name = common.StringOrNull(inst.Name)
	model.Cloud = common.StringOrNull(inst.Cloud)
	model.Region = common.StringOrNull(inst.Region)
	model.ShadeInstanceType = common.StringOrNull(inst.InstanceType)
	model.ShadeCloud = types.BoolValue(inst.ShadeCloud)
	model.OS = common.StringOrNull(inst.OS)
	model.TemplateID = common.StringOrNull(inst.TemplateID)
	model.SSHKeyID = common.StringOrNull(inst.SSHKeyID)
	model.AutoDelete = types.BoolValue(inst.AutoDelete)
	model.Alert = types.BoolValue(inst.Alert)
	model.Tags = common.SetStringsOrEmpty(inst.Tags, model.Tags)
	model.Envs = common.MapStringsOrEmpty(inst.Envs, model.Envs)
	if len(inst.VolumeIDs) > 0 {
		model.VolumeIDs = common.SetStrings(inst.VolumeIDs)
	}
	model.Status = common.StringOrNull(inst.Status)
	model.IP = common.StringOrNull(inst.IP)
	return true
}
