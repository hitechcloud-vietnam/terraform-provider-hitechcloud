// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vstorage

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
	_ resource.Resource                = (*s3SubuserResource)(nil)
	_ resource.ResourceWithConfigure   = (*s3SubuserResource)(nil)
	_ resource.ResourceWithImportState = (*s3SubuserResource)(nil)
)

// S3SubuserResource returns the hitechcloud_s3_subuser resource.
func S3SubuserResource() resource.Resource {
	return &s3SubuserResource{}
}

type s3SubuserResource struct {
	client *client.Client
}

type s3SubuserModel struct {
	ID        types.String `tfsdk:"id"`
	ServiceID types.String `tfsdk:"service_id"`
	Name      types.String `tfsdk:"name"`
	Access    types.String `tfsdk:"access"`
	AccessKey types.String `tfsdk:"access_key"`
	SecretKey types.String `tfsdk:"secret_key"`
}

func (r *s3SubuserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_subuser"
}

func (r *s3SubuserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a sub-user of the HiTechCloud Ceph S3 service " +
			"(`POST /api/service/{service_id}/s3/subusers`). The secret key is " +
			"returned by the API only once at creation; it is stored in state and " +
			"will read as `null` for imported sub-users. The API has no update " +
			"endpoint, so changes force replacement.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `<service_id>/<name>`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the S3 service. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Sub-user name. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"access": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Permission level: `read`, `write`, `readwrite` or `full`. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"access_key": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Access key id of the sub-user.",
			},
			"secret_key": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Secret access key (only returned at creation; null for imported sub-users).",
			},
		},
	}
}

func (r *s3SubuserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *s3SubuserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan s3SubuserModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	access := plan.Access.ValueString()
	if access == "" {
		access = "readwrite"
	}

	sub, err := r.client.CreateS3Subuser(ctx, plan.ServiceID.ValueString(), plan.Name.ValueString(), access)
	if err != nil {
		resp.Diagnostics.AddError("Error creating S3 sub-user", err.Error())
		return
	}

	plan.ID = types.StringValue(client.FormatID(plan.ServiceID.ValueString(), sub.Name))
	plan.Access = common.StringOrNull(sub.Access)
	plan.AccessKey = common.StringOrNull(sub.AccessKey)
	plan.SecretKey = common.StringOrNull(sub.SecretKey)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *s3SubuserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state s3SubuserModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, name, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	sub, err := r.client.GetS3Subuser(ctx, serviceID, name)
	if client.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading S3 sub-user", err.Error())
		return
	}

	state.ServiceID = types.StringValue(serviceID)
	state.Name = types.StringValue(sub.Name)
	state.Access = common.StringOrNull(sub.Access)
	state.AccessKey = common.StringOrNull(sub.AccessKey)
	// The secret key is never returned after creation; keep whatever state has.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is unreachable: every attribute forces replacement.
func (r *s3SubuserResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"The HiTechCloud API provides no sub-user update endpoint; the sub-user must be replaced instead.",
	)
}

func (r *s3SubuserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state s3SubuserModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, name, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	if err := r.client.DeleteS3Subuser(ctx, serviceID, name); client.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error deleting S3 sub-user", err.Error())
	}
}

func (r *s3SubuserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := client.SplitID(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[1])...)
}

func (r *s3SubuserResource) parseID(_ context.Context, state *s3SubuserModel, diags *diag.Diagnostics) (string, string, bool) {
	if state.ServiceID.ValueString() != "" && state.Name.ValueString() != "" {
		return state.ServiceID.ValueString(), state.Name.ValueString(), true
	}
	parts, err := client.SplitID(state.ID.ValueString(), 2)
	if err != nil {
		diags.AddError("Unexpected Resource ID", err.Error())
		return "", "", false
	}
	return parts[0], parts[1], true
}
