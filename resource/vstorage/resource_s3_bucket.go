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
)

var (
	_ resource.Resource                = (*s3BucketResource)(nil)
	_ resource.ResourceWithConfigure   = (*s3BucketResource)(nil)
	_ resource.ResourceWithImportState = (*s3BucketResource)(nil)
)

// S3BucketResource returns the hitechcloud_s3_bucket resource.
func S3BucketResource() resource.Resource {
	return &s3BucketResource{}
}

type s3BucketResource struct {
	client *client.Client
}

type s3BucketModel struct {
	ID        types.String `tfsdk:"id"`
	ServiceID types.String `tfsdk:"service_id"`
	Name      types.String `tfsdk:"name"`
	Bucket    types.String `tfsdk:"bucket"`
	Purge     types.Bool   `tfsdk:"purge"`
}

func (r *s3BucketResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_bucket"
}

func (r *s3BucketResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a bucket of the HiTechCloud Ceph S3 service " +
			"(`POST /api/service/{service_id}/s3/buckets`). Bucket names are " +
			"normalised by the API (lower case, optional service prefix), so the " +
			"effective name is exported as `bucket`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `<service_id>/<bucket>`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the S3 service. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Requested bucket name. Forces replacement (buckets cannot be renamed).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"bucket": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Effective bucket name as reported by the API (may include a service prefix).",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"purge": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "When `true`, destroying the resource also deletes the objects in the bucket. Defaults to `false` (non-empty buckets are rejected by the API).",
			},
		},
	}
}

func (r *s3BucketResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *s3BucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan s3BucketModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucket, err := r.client.CreateS3Bucket(ctx, plan.ServiceID.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating S3 bucket", err.Error())
		return
	}

	plan.Bucket = types.StringValue(bucket)
	plan.ID = types.StringValue(client.FormatID(plan.ServiceID.ValueString(), bucket))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *s3BucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state s3BucketModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, bucket, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	found, err := r.client.GetS3Bucket(ctx, serviceID, bucket)
	if client.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading S3 bucket", err.Error())
		return
	}

	state.ServiceID = types.StringValue(serviceID)
	state.Bucket = types.StringValue(found.Name)
	state.ID = types.StringValue(client.FormatID(serviceID, found.Name))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is unreachable: every attribute forces replacement.
func (r *s3BucketResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"S3 buckets cannot be renamed through the HiTechCloud API; the bucket must be replaced instead.",
	)
}

func (r *s3BucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state s3BucketModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, bucket, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	purge := !state.Purge.IsNull() && !state.Purge.IsUnknown() && state.Purge.ValueBool()
	if err := r.client.DeleteS3Bucket(ctx, serviceID, bucket, purge); client.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error deleting S3 bucket", err.Error())
	}
}

func (r *s3BucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := client.SplitID(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("bucket"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[1])...)
}

func (r *s3BucketResource) parseID(_ context.Context, state *s3BucketModel, diags *diag.Diagnostics) (string, string, bool) {
	if state.ServiceID.ValueString() != "" && state.Bucket.ValueString() != "" {
		return state.ServiceID.ValueString(), state.Bucket.ValueString(), true
	}
	parts, err := client.SplitID(state.ID.ValueString(), 2)
	if err != nil {
		diags.AddError("Unexpected Resource ID", err.Error())
		return "", "", false
	}
	return parts[0], parts[1], true
}
