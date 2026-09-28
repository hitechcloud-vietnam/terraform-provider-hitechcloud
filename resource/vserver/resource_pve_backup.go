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
	_ resource.Resource                = (*pveBackupResource)(nil)
	_ resource.ResourceWithConfigure   = (*pveBackupResource)(nil)
	_ resource.ResourceWithImportState = (*pveBackupResource)(nil)
)

// PVEBackupResource returns the hitechcloud_pve_backup resource.
func PVEBackupResource() resource.Resource {
	return &pveBackupResource{}
}

type pveBackupResource struct {
	client *client.Client
}

type pveBackupModel struct {
	ID        types.String `tfsdk:"id"`
	ServiceID types.String `tfsdk:"service_id"`
	Mode      types.String `tfsdk:"mode"`
	Notes     types.String `tfsdk:"notes"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func (r *pveBackupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pve_backup"
}

func (r *pveBackupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Takes a backup of a HiTechCloud Proxmox service " +
			"(`POST /api/service/{id}/htcpve/backups`). Backups are immutable " +
			"records: changing any attribute forces a replacement. Deleting the " +
			"resource removes it from Terraform state only (the API has no backup " +
			"deletion endpoint).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Backup identifier returned by the API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HiTechCloud service ID of the Proxmox service. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mode": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Backup mode, for example `snapshot`. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"notes": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Free-text notes for the backup. Changing this forces a new resource.",
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

func (r *pveBackupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *pveBackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan pveBackupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	backup, err := r.client.CreatePVEBackup(ctx,
		plan.ServiceID.ValueString(),
		plan.Mode.ValueString(),
		plan.Notes.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating PVE backup", err.Error())
		return
	}

	id := client.FirstString(backup, "id", "backup", "uid")
	if id == "" {
		id = client.FirstString(backup, "created_at", "created", "time")
	}
	plan.ID = types.StringValue(id)
	plan.CreatedAt = common.StringOrNull(client.FirstString(backup, "created_at", "created", "time"))
	if v := client.FirstString(backup, "mode", "type"); v != "" {
		plan.Mode = types.StringValue(v)
	}
	if v := client.FirstString(backup, "notes"); v != "" {
		plan.Notes = types.StringValue(v)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *pveBackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pveBackupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	backups, err := r.client.ListPVEBackups(ctx, state.ServiceID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading PVE backups", err.Error())
		return
	}

	m, ok := client.FindByID(backups, state.ID.ValueString(), "id", "backup", "uid")
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}
	state.Mode = common.StringOrNull(client.FirstString(m, "mode", "type"))
	state.Notes = common.StringOrNull(client.FirstString(m, "notes", "description"))
	state.CreatedAt = common.StringOrNull(client.FirstString(m, "created_at", "created", "time"))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Backups are immutable: Update is only reached through forced replacement.
func (r *pveBackupResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"PVE backups are immutable",
		"The API does not support editing backups; this resource only supports create and read.",
	)
}

func (r *pveBackupResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Backup not deleted",
		"The HiTechCloud API does not expose a backup deletion endpoint; the backup was removed from Terraform state only.",
	)
}

func (r *pveBackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	parts, err := client.SplitID(id, 2)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf(
			"Expected import ID in the format service_id/backup_id: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), parts[0])...)
}
