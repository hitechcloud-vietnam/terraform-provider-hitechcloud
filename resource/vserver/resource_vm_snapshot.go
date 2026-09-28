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
	_ resource.Resource                = (*vmSnapshotResource)(nil)
	_ resource.ResourceWithConfigure   = (*vmSnapshotResource)(nil)
	_ resource.ResourceWithImportState = (*vmSnapshotResource)(nil)
)

// VMSnapshotResource returns the hitechcloud_vm_snapshot resource.
func VMSnapshotResource() resource.Resource {
	return &vmSnapshotResource{}
}

type vmSnapshotResource struct {
	client *client.Client
}

type vmSnapshotModel struct {
	ID          types.String `tfsdk:"id"`
	ServiceID   types.String `tfsdk:"service_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

func (r *vmSnapshotResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm_snapshot"
}

func (r *vmSnapshotResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Takes a snapshot of a HiTechCloud Proxmox virtual machine " +
			"(`POST /api/service/{id}/htcpve/snapshots`). Snapshots are immutable " +
			"records: changing any attribute forces a replacement.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Snapshot identifier returned by the API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HiTechCloud service ID of the Proxmox VM. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Snapshot name; invalid characters are replaced by the API. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Free-text description of the snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Creation timestamp reported by the API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *vmSnapshotResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *vmSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vmSnapshotModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	snap, err := r.client.CreatePVESnapshot(ctx,
		plan.ServiceID.ValueString(),
		plan.Name.ValueString(),
		plan.Description.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VM snapshot", err.Error())
		return
	}

	id := client.FirstString(snap, "id", "snapshot", "name")
	if id == "" {
		id = plan.Name.ValueString()
	}
	plan.ID = types.StringValue(id)
	plan.CreatedAt = common.StringOrNull(client.FirstString(snap, "created_at", "created", "time"))
	if v := client.FirstString(snap, "description"); v != "" {
		plan.Description = types.StringValue(v)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *vmSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vmSnapshotModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	snaps, err := r.client.ListPVESnapshots(ctx, state.ServiceID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading VM snapshots", err.Error())
		return
	}

	m, ok := client.FindByID(snaps, state.ID.ValueString(), "id", "snapshot", "name")
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}
	state.Name = types.StringValue(client.FirstString(m, "name", "snapshot", "id"))
	state.Description = common.StringOrNull(client.FirstString(m, "description", "desc"))
	state.CreatedAt = common.StringOrNull(client.FirstString(m, "created_at", "created", "time"))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Snapshots are immutable: Update is only reached through forced replacement.
func (r *vmSnapshotResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"VM snapshots are immutable",
		"The API does not support editing snapshots; this resource only supports create, read and delete.",
	)
}

func (r *vmSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vmSnapshotModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The API has no snapshot delete endpoint exposed; removing from state is
	// the supported behaviour and matches provider convention for action-style
	// resources.
	resp.Diagnostics.AddWarning(
		"Snapshot not deleted",
		"The HiTechCloud API does not expose a snapshot deletion endpoint; the snapshot was removed from Terraform state only.",
	)
}

func (r *vmSnapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	parts, err := client.SplitID(id, 2)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf(
			"Expected import ID in the format service_id/snapshot_id: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), parts[0])...)
}
