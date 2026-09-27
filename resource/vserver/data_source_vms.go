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
	_ datasource.DataSource              = (*vmsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*vmsDataSource)(nil)
)

// VMsDataSource returns the hitechcloud_vms data source.
func VMsDataSource() datasource.DataSource {
	return &vmsDataSource{}
}

type vmsDataSource struct {
	client *client.Client
}

type vmsModel struct {
	ID        types.String `tfsdk:"id"`
	ServiceID types.String `tfsdk:"service_id"`
	VMs       []vmItem     `tfsdk:"vms"`
}

type vmItem struct {
	ID         types.String `tfsdk:"id"`
	Label      types.String `tfsdk:"label"`
	Status     types.String `tfsdk:"status"`
	Hostname   types.String `tfsdk:"hostname"`
	TemplateID types.String `tfsdk:"template_id"`
	Memory     types.Int64  `tfsdk:"memory"`
	CPU        types.Int64  `tfsdk:"cpu"`
	Disk       types.Int64  `tfsdk:"disk"`
	IPs        []string     `tfsdk:"ips"`
}

func (d *vmsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vms"
}

func (d *vmsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the virtual machines of a HiTechCloud cloud service " +
			"(`GET /api/service/{service_id}/vms`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (same as `service_id`).",
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the cloud service.",
			},
			"vms": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Virtual machines of the service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.StringAttribute{Computed: true, MarkdownDescription: "VM identifier."},
						"label":       schema.StringAttribute{Computed: true, MarkdownDescription: "VM label."},
						"status":      schema.StringAttribute{Computed: true, MarkdownDescription: "VM status."},
						"hostname":    schema.StringAttribute{Computed: true, MarkdownDescription: "VM hostname."},
						"template_id": schema.StringAttribute{Computed: true, MarkdownDescription: "OS template id."},
						"memory":      schema.Int64Attribute{Computed: true, MarkdownDescription: "Memory in MB."},
						"cpu":         schema.Int64Attribute{Computed: true, MarkdownDescription: "CPU cores."},
						"disk":        schema.Int64Attribute{Computed: true, MarkdownDescription: "Disk size in GB."},
						"ips": schema.SetAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "IP addresses.",
						},
					},
				},
			},
		},
	}
}

func (d *vmsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vmsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data vmsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := data.ServiceID.ValueString()
	vms, err := d.client.ListVMs(ctx, serviceID)
	if err != nil {
		resp.Diagnostics.AddError("Error listing VMs", err.Error())
		return
	}

	items := make([]vmItem, 0, len(vms))
	for _, vm := range vms {
		items = append(items, vmItem{
			ID:         common.StringOrNull(vm.ID),
			Label:      common.StringOrNull(vm.Label),
			Status:     common.StringOrNull(vm.Status),
			Hostname:   common.StringOrNull(vm.Hostname),
			TemplateID: common.StringOrNull(vm.TemplateID),
			Memory:     common.Int64OrNull(vm.Memory),
			CPU:        common.Int64OrNull(vm.CPU),
			Disk:       common.Int64OrNull(vm.Disk),
			IPs:        vm.IPs,
		})
	}

	data.ID = types.StringValue(serviceID)
	data.VMs = items
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
