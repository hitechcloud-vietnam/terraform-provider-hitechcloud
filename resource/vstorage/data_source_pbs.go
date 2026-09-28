// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vstorage

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
	_ datasource.DataSource              = (*pbsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*pbsDataSource)(nil)
)

// PBSDataSource returns the hitechcloud_pbs data source.
func PBSDataSource() datasource.DataSource {
	return &pbsDataSource{}
}

type pbsDataSource struct {
	client *client.Client
}

type pbsModel struct {
	ID        types.String   `tfsdk:"id"`
	ServiceID types.String   `tfsdk:"service_id"`
	Endpoint  types.String   `tfsdk:"endpoint"`
	Username  types.String   `tfsdk:"username"`
	Namespace types.String   `tfsdk:"namespace"`
	Usage     types.Map      `tfsdk:"usage"`
	Metrics   types.Map      `tfsdk:"metrics"`
	Snapshots []pbsSnapshot  `tfsdk:"snapshots"`
	Groups    []types.String `tfsdk:"groups"`
}

type pbsSnapshot struct {
	ID        types.String `tfsdk:"id"`
	Group     types.String `tfsdk:"group"`
	BackupID  types.String `tfsdk:"backup_id"`
	Size      types.Int64  `tfsdk:"size"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func (d *pbsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pbs"
}

func (d *pbsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads a Proxmox Backup Server (PBS) service: connection info, " +
			"credentials, usage metrics and snapshots " +
			"(`GET /api/service/{id}/pbs`, `/pbs/credentials`, `/pbs/usage`, `/pbs/metrics`, " +
			"`/pbs/snapshots`, `/pbs/groups`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The service ID (same as `service_id`).",
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HiTechCloud service ID (`hb_accounts.id`) of the PBS service.",
			},
			"endpoint":  schema.StringAttribute{Computed: true, MarkdownDescription: "PBS endpoint / hostname."},
			"username":  schema.StringAttribute{Computed: true, MarkdownDescription: "PBS access username."},
			"namespace": schema.StringAttribute{Computed: true, MarkdownDescription: "PBS namespace of the customer."},
			"usage":     schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Billing-relevant usage counters (`backup_space`, `snapshots`, `backup_groups`) with units."},
			"metrics":   schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Raw metric values keyed by metric name."},
			"snapshots": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Backups in the customer's namespace (newest first).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.StringAttribute{Computed: true, MarkdownDescription: "Snapshot identifier."},
						"group":      schema.StringAttribute{Computed: true, MarkdownDescription: "Backup group."},
						"backup_id":  schema.StringAttribute{Computed: true, MarkdownDescription: "Backup ID inside the group."},
						"size":       schema.Int64Attribute{Computed: true, MarkdownDescription: "Snapshot size in bytes; null when PBS does not report a size."},
						"created_at": schema.StringAttribute{Computed: true, MarkdownDescription: "Backup timestamp."},
					},
				},
			},
			"groups": schema.ListAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "Backup groups present in the namespace.",
			},
		},
	}
}

func (d *pbsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *pbsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data pbsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := data.ServiceID.ValueString()
	data.ID = types.StringValue(serviceID)

	if info, err := d.client.GetPBSConnectionInfo(ctx, serviceID); err == nil {
		data.Endpoint = common.StringOrNull(client.FirstString(info, "endpoint", "host", "hostname", "url"))
		data.Username = common.StringOrNull(client.FirstString(info, "username", "user"))
		data.Namespace = common.StringOrNull(client.FirstString(info, "namespace", "ns"))
	}
	if creds, err := d.client.GetPBSCredentials(ctx, serviceID); err == nil {
		if data.Username.IsNull() {
			data.Username = common.StringOrNull(client.FirstString(creds, "username", "user"))
		}
	}
	if usage, err := d.client.GetPBSUsage(ctx, serviceID); err == nil {
		data.Usage = common.MapStrings(client.StringMap(usage))
	}
	if metrics, err := d.client.GetPBSMetrics(ctx, serviceID, ""); err == nil {
		data.Metrics = common.MapStrings(client.StringMap(metrics))
	}

	snaps, err := d.client.ListPBSSnapshots(ctx, serviceID, "")
	if err != nil {
		resp.Diagnostics.AddError("Error reading PBS snapshots", err.Error())
		return
	}
	data.Snapshots = make([]pbsSnapshot, 0, len(snaps))
	for _, item := range snaps {
		m, ok := client.AsMap(item)
		if !ok {
			continue
		}
		entry := pbsSnapshot{
			ID:        common.StringOrNull(client.FirstString(m, "id", "snapshot", "uid")),
			Group:     common.StringOrNull(client.FirstString(m, "group", "backup_group")),
			BackupID:  common.StringOrNull(client.FirstString(m, "backup_id", "backup-id", "vmid")),
			CreatedAt: common.StringOrNull(client.FirstString(m, "created_at", "created", "time", "ctime")),
		}
		if s := client.FirstInt64(m, "size"); s != 0 {
			entry.Size = types.Int64Value(s)
		} else if raw, exists := m["size"]; exists && raw == nil {
			entry.Size = types.Int64Null()
		}
		data.Snapshots = append(data.Snapshots, entry)
	}

	groups := []string{}
	if list, err := d.client.ListPBSGroups(ctx, serviceID); err == nil {
		for _, g := range list {
			if s := client.FirstString(g, "group", "name", "id"); s != "" {
				groups = append(groups, s)
			}
		}
	}
	data.Groups = make([]types.String, 0, len(groups))
	for _, g := range groups {
		data.Groups = append(data.Groups, types.StringValue(g))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
