// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/client"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/common"
)

var (
	_ datasource.DataSource              = (*pveDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*pveDataSource)(nil)
)

// PVEDataSource returns the hitechcloud_pve data source.
func PVEDataSource() datasource.DataSource {
	return &pveDataSource{}
}

type pveDataSource struct {
	client *client.Client
}

type pveModel struct {
	ID        types.String  `tfsdk:"id"`
	ServiceID types.String  `tfsdk:"service_id"`
	Status    types.Map     `tfsdk:"status"`
	Usage     types.Map     `tfsdk:"usage"`
	VMs       []pveVM       `tfsdk:"vms"`
	IPs       []pveIP       `tfsdk:"ips"`
	Backups   []pveBackup   `tfsdk:"backups"`
	Snapshots []pveSnapshot `tfsdk:"snapshots"`
}

type pveVM struct {
	VMID     types.String `tfsdk:"vmid"`
	Name     types.String `tfsdk:"name"`
	Status   types.String `tfsdk:"status"`
	CPUs     types.Int64  `tfsdk:"cpus"`
	MemoryMB types.Int64  `tfsdk:"memory_mb"`
	DiskGB   types.Int64  `tfsdk:"disk_gb"`
}

type pveIP struct {
	IP       types.String `tfsdk:"ip"`
	Hostname types.String `tfsdk:"hostname"`
	Type     types.String `tfsdk:"type"`
}

type pveBackup struct {
	ID        types.String `tfsdk:"id"`
	VMID      types.String `tfsdk:"vmid"`
	Mode      types.String `tfsdk:"mode"`
	Size      types.Int64  `tfsdk:"size"`
	Notes     types.String `tfsdk:"notes"`
	CreatedAt types.String `tfsdk:"created_at"`
}

type pveSnapshot struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

func (d *pveDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pve"
}

func (d *pveDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads a HiTechCloud Proxmox (PVE) service: status, VM inventory, " +
			"IPs, backups, snapshots and usage " +
			"(`GET /api/service/{id}/htcpve/status`, `/htcpve/vms`, `/htcpve/ips`, " +
			"`/htcpve/backups`, `/htcpve/snapshots`, `/htcpve/usage`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The service ID (same as `service_id`).",
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HiTechCloud service ID of the Proxmox service.",
			},
			"status": schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Node-level status values."},
			"usage":  schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Billing-relevant usage counters with units."},
			"vms": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Virtual machines of the service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vmid":      schema.StringAttribute{Computed: true, MarkdownDescription: "Proxmox VM ID."},
						"name":      schema.StringAttribute{Computed: true, MarkdownDescription: "VM name."},
						"status":    schema.StringAttribute{Computed: true, MarkdownDescription: "VM status (running/stopped/...)."},
						"cpus":      schema.Int64Attribute{Computed: true, MarkdownDescription: "Number of vCPUs."},
						"memory_mb": schema.Int64Attribute{Computed: true, MarkdownDescription: "Memory in MB."},
						"disk_gb":   schema.Int64Attribute{Computed: true, MarkdownDescription: "Disk size in GB."},
					},
				},
			},
			"ips": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "IP addresses assigned to the service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip":       schema.StringAttribute{Computed: true, MarkdownDescription: "IP address."},
						"hostname": schema.StringAttribute{Computed: true, MarkdownDescription: "Reverse DNS hostname."},
						"type":     schema.StringAttribute{Computed: true, MarkdownDescription: "Address type (ipv4/ipv6)."},
					},
				},
			},
			"backups": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Backups of the service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.StringAttribute{Computed: true, MarkdownDescription: "Backup identifier."},
						"vmid":       schema.StringAttribute{Computed: true, MarkdownDescription: "Backed-up VM ID."},
						"mode":       schema.StringAttribute{Computed: true, MarkdownDescription: "Backup mode (snapshot/...)."},
						"size":       schema.Int64Attribute{Computed: true, MarkdownDescription: "Backup size in bytes."},
						"notes":      schema.StringAttribute{Computed: true, MarkdownDescription: "Backup notes."},
						"created_at": schema.StringAttribute{Computed: true, MarkdownDescription: "Backup timestamp."},
					},
				},
			},
			"snapshots": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Snapshots of the machine, excluding the synthetic \"current\" entry.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true, MarkdownDescription: "Snapshot identifier."},
						"name":        schema.StringAttribute{Computed: true, MarkdownDescription: "Snapshot name."},
						"description": schema.StringAttribute{Computed: true, MarkdownDescription: "Snapshot description."},
						"created_at":  schema.StringAttribute{Computed: true, MarkdownDescription: "Snapshot timestamp."},
					},
				},
			},
		},
	}
}

func (d *pveDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	cli, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got %T", req.ProviderData),
		)
		return
	}
	d.client = cli
}

func (d *pveDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data pveModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := data.ServiceID.ValueString()
	data.ID = types.StringValue(serviceID)

	if status, err := d.client.GetPVEStatus(ctx, serviceID); err == nil {
		data.Status = common.MapStrings(client.StringMap(status))
	}
	if usage, err := d.client.GetPVEUsage(ctx, serviceID); err == nil {
		data.Usage = common.MapStrings(client.StringMap(usage))
	}

	if vms, err := d.client.ListPVEVMs(ctx, serviceID); err == nil {
		data.VMs = make([]pveVM, 0, len(vms))
		for _, item := range vms {
			m, ok := client.AsMap(item)
			if !ok {
				continue
			}
			data.VMs = append(data.VMs, pveVM{
				VMID:     types.StringValue(client.FirstString(m, "vmid", "vm_id", "id")),
				Name:     common.StringOrNull(client.FirstString(m, "name", "hostname")),
				Status:   common.StringOrNull(client.FirstString(m, "status", "state")),
				CPUs:     types.Int64Value(client.FirstInt64(m, "cpus", "cpu", "cores")),
				MemoryMB: types.Int64Value(client.FirstInt64(m, "memory_mb", "memory", "mem")),
				DiskGB:   types.Int64Value(client.FirstInt64(m, "disk_gb", "disk", "storage")),
			})
		}
	}

	if ips, err := d.client.ListPVEIPs(ctx, serviceID); err == nil {
		data.IPs = make([]pveIP, 0, len(ips))
		for _, item := range ips {
			m, ok := client.AsMap(item)
			if !ok {
				continue
			}
			data.IPs = append(data.IPs, pveIP{
				IP:       types.StringValue(client.FirstString(m, "ip", "address", "addr")),
				Hostname: common.StringOrNull(client.FirstString(m, "hostname", "host", "ptr")),
				Type:     common.StringOrNull(client.FirstString(m, "type", "family", "version")),
			})
		}
	}

	if backups, err := d.client.ListPVEBackups(ctx, serviceID); err == nil {
		data.Backups = make([]pveBackup, 0, len(backups))
		for _, item := range backups {
			m, ok := client.AsMap(item)
			if !ok {
				continue
			}
			data.Backups = append(data.Backups, pveBackup{
				ID:        common.StringOrNull(client.FirstString(m, "id", "backup", "uid")),
				VMID:      common.StringOrNull(client.FirstString(m, "vmid", "vm_id")),
				Mode:      common.StringOrNull(client.FirstString(m, "mode", "type")),
				Size:      types.Int64Value(client.FirstInt64(m, "size")),
				Notes:     common.StringOrNull(client.FirstString(m, "notes", "description")),
				CreatedAt: common.StringOrNull(client.FirstString(m, "created_at", "created", "time")),
			})
		}
	}

	snaps, err := d.client.ListPVESnapshots(ctx, serviceID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading PVE snapshots", err.Error())
		return
	}
	data.Snapshots = make([]pveSnapshot, 0, len(snaps))
	for _, item := range snaps {
		m, ok := client.AsMap(item)
		if !ok {
			continue
		}
		data.Snapshots = append(data.Snapshots, pveSnapshot{
			ID:          common.StringOrNull(client.FirstString(m, "id", "snapshot", "name")),
			Name:        common.StringOrNull(client.FirstString(m, "name", "snapshot")),
			Description: common.StringOrNull(client.FirstString(m, "description", "desc")),
			CreatedAt:   common.StringOrNull(client.FirstString(m, "created_at", "created", "time")),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
