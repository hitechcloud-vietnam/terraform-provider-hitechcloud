// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vaifactory

import (
	"context"
	"fmt"

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
	_ resource.Resource                = (*aiVolumeResource)(nil)
	_ resource.ResourceWithConfigure   = (*aiVolumeResource)(nil)
	_ resource.ResourceWithImportState = (*aiVolumeResource)(nil)
)

// AIVolumeResource returns the hitechcloud_ai_volume resource.
func AIVolumeResource() resource.Resource {
	return &aiVolumeResource{}
}

type aiVolumeResource struct {
	client *client.Client
}

type aiVolumeModel struct {
	ID        types.String `tfsdk:"id"`
	ServiceID types.String `tfsdk:"service_id"`
	VolumeID  types.String `tfsdk:"volume_id"`
	Name      types.String `tfsdk:"name"`
	Cloud     types.String `tfsdk:"cloud"`
	Region    types.String `tfsdk:"region"`
	SizeInGB  types.Int64  `tfsdk:"size_in_gb"`
	Status    types.String `tfsdk:"status"`
}

func (r *aiVolumeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ai_volume"
}

func (r *aiVolumeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a HiTechCloud AI Factory storage volume " +
			"(`POST /api/service/{service_id}/volumes`). The API provides no volume " +
			"update endpoint, so every change forces replacement.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `<service_id>/<volume_id>`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the AI Factory service. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"volume_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Volume identifier.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Volume name. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"cloud": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Cloud provider. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"region": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Region. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"size_in_gb": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Volume size in GB. Forces replacement.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current volume status.",
			},
		},
	}
}

func (r *aiVolumeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *aiVolumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aiVolumeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := client.AIVolume{
		Name:     plan.Name.ValueString(),
		Cloud:    plan.Cloud.ValueString(),
		Region:   plan.Region.ValueString(),
		SizeInGB: plan.SizeInGB.ValueInt64(),
	}
	volumeID, err := r.client.CreateAIVolume(ctx, plan.ServiceID.ValueString(), in)
	if err != nil {
		resp.Diagnostics.AddError("Error creating volume", err.Error())
		return
	}

	plan.VolumeID = types.StringValue(volumeID)
	plan.ID = types.StringValue(client.FormatID(plan.ServiceID.ValueString(), volumeID))
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *aiVolumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state aiVolumeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, volumeID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}
	state.ServiceID = types.StringValue(serviceID)
	state.VolumeID = types.StringValue(volumeID)
	if !r.refresh(ctx, &state, &resp.Diagnostics) {
		resp.State.RemoveResource(ctx)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is unreachable: every attribute forces replacement.
func (r *aiVolumeResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"The HiTechCloud API provides no volume update endpoint; the volume must be replaced instead.",
	)
}

func (r *aiVolumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state aiVolumeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, volumeID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	if err := r.client.DeleteAIVolume(ctx, serviceID, volumeID); client.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error deleting volume", err.Error())
	}
}

func (r *aiVolumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := client.SplitID(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("volume_id"), parts[1])...)
}

func (r *aiVolumeResource) parseID(_ context.Context, state *aiVolumeModel, diags *diag.Diagnostics) (string, string, bool) {
	if state.ServiceID.ValueString() != "" && state.VolumeID.ValueString() != "" {
		return state.ServiceID.ValueString(), state.VolumeID.ValueString(), true
	}
	parts, err := client.SplitID(state.ID.ValueString(), 2)
	if err != nil {
		diags.AddError("Unexpected Resource ID", err.Error())
		return "", "", false
	}
	return parts[0], parts[1], true
}

// refresh reloads the volume; returns false when it no longer exists.
func (r *aiVolumeResource) refresh(ctx context.Context, model *aiVolumeModel, diags *diag.Diagnostics) bool {
	vol, err := r.client.GetAIVolume(ctx, model.ServiceID.ValueString(), model.VolumeID.ValueString())
	if client.IsNotFound(err) {
		return false
	}
	if err != nil {
		diags.AddError("Error reading volume", err.Error())
		return true
	}
	model.Name = common.StringOrNull(vol.Name)
	model.Cloud = common.StringOrNull(vol.Cloud)
	model.Region = common.StringOrNull(vol.Region)
	model.SizeInGB = common.Int64OrNull(vol.SizeInGB)
	model.Status = common.StringOrNull(vol.Status)
	return true
}
