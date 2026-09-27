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
	_ resource.Resource                = (*aiClusterResource)(nil)
	_ resource.ResourceWithConfigure   = (*aiClusterResource)(nil)
	_ resource.ResourceWithImportState = (*aiClusterResource)(nil)
)

// AIClusterResource returns the hitechcloud_ai_cluster resource.
func AIClusterResource() resource.Resource {
	return &aiClusterResource{}
}

type aiClusterResource struct {
	client *client.Client
}

type aiClusterModel struct {
	ID           types.String `tfsdk:"id"`
	ServiceID    types.String `tfsdk:"service_id"`
	ClusterID    types.String `tfsdk:"cluster_id"`
	Name         types.String `tfsdk:"name"`
	Cloud        types.String `tfsdk:"cloud"`
	Region       types.String `tfsdk:"region"`
	ClusterType  types.String `tfsdk:"cluster_type"`
	NumInstances types.Int64  `tfsdk:"num_instances"`
	SSHKeyID     types.String `tfsdk:"ssh_key_id"`
	OS           types.String `tfsdk:"os"`
	Status       types.String `tfsdk:"status"`
}

func (r *aiClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ai_cluster"
}

func (r *aiClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a HiTechCloud AI Factory GPU cluster " +
			"(`POST /api/service/{service_id}/clusters`). The API provides no cluster " +
			"update endpoint, so every change forces replacement.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `<service_id>/<cluster_id>`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the AI Factory service. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"cluster_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Cluster identifier.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Cluster name. Forces replacement.",
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
			"cluster_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Cluster type (see `hitechcloud_ai_cluster_types`). Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"num_instances": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Number of instances in the cluster. Forces replacement.",
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"ssh_key_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "SSH key id to inject. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"os": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Operating system / image. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current cluster status.",
			},
		},
	}
}

func (r *aiClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *aiClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aiClusterModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := client.AICluster{
		Name:         plan.Name.ValueString(),
		Cloud:        plan.Cloud.ValueString(),
		Region:       plan.Region.ValueString(),
		ClusterType:  plan.ClusterType.ValueString(),
		NumInstances: plan.NumInstances.ValueInt64(),
		SSHKeyID:     plan.SSHKeyID.ValueString(),
		OS:           plan.OS.ValueString(),
	}
	clusterID, err := r.client.CreateAICluster(ctx, plan.ServiceID.ValueString(), in)
	if err != nil {
		resp.Diagnostics.AddError("Error creating cluster", err.Error())
		return
	}

	plan.ClusterID = types.StringValue(clusterID)
	plan.ID = types.StringValue(client.FormatID(plan.ServiceID.ValueString(), clusterID))
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *aiClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state aiClusterModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, clusterID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}
	state.ServiceID = types.StringValue(serviceID)
	state.ClusterID = types.StringValue(clusterID)
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
func (r *aiClusterResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"The HiTechCloud API provides no cluster update endpoint; the cluster must be replaced instead.",
	)
}

func (r *aiClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state aiClusterModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, clusterID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	if err := r.client.DeleteAICluster(ctx, serviceID, clusterID); client.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error deleting cluster", err.Error())
	}
}

func (r *aiClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := client.SplitID(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cluster_id"), parts[1])...)
}

func (r *aiClusterResource) parseID(_ context.Context, state *aiClusterModel, diags *diag.Diagnostics) (string, string, bool) {
	if state.ServiceID.ValueString() != "" && state.ClusterID.ValueString() != "" {
		return state.ServiceID.ValueString(), state.ClusterID.ValueString(), true
	}
	parts, err := client.SplitID(state.ID.ValueString(), 2)
	if err != nil {
		diags.AddError("Unexpected Resource ID", err.Error())
		return "", "", false
	}
	return parts[0], parts[1], true
}

// refresh reloads the cluster; returns false when it no longer exists.
func (r *aiClusterResource) refresh(ctx context.Context, model *aiClusterModel, diags *diag.Diagnostics) bool {
	cl, err := r.client.GetAICluster(ctx, model.ServiceID.ValueString(), model.ClusterID.ValueString())
	if client.IsNotFound(err) {
		return false
	}
	if err != nil {
		diags.AddError("Error reading cluster", err.Error())
		return true
	}
	model.Name = common.StringOrNull(cl.Name)
	model.Cloud = common.StringOrNull(cl.Cloud)
	model.Region = common.StringOrNull(cl.Region)
	model.ClusterType = common.StringOrNull(cl.ClusterType)
	model.NumInstances = common.Int64OrNull(cl.NumInstances)
	model.SSHKeyID = common.StringOrNull(cl.SSHKeyID)
	model.OS = common.StringOrNull(cl.OS)
	model.Status = common.StringOrNull(cl.Status)
	return true
}
