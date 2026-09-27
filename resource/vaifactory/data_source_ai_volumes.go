// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vaifactory

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
	_ datasource.DataSource              = (*aiVolumesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*aiVolumesDataSource)(nil)
)

// AIVolumesDataSource returns the hitechcloud_ai_volumes data source.
func AIVolumesDataSource() datasource.DataSource {
	return &aiVolumesDataSource{}
}

type aiVolumesDataSource struct {
	client *client.Client
}

type aiVolumesModel struct {
	ID        types.String   `tfsdk:"id"`
	ServiceID types.String   `tfsdk:"service_id"`
	Volumes   []aiVolumeItem `tfsdk:"volumes"`
}

type aiVolumeItem struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Cloud    types.String `tfsdk:"cloud"`
	Region   types.String `tfsdk:"region"`
	SizeInGB types.Int64  `tfsdk:"size_in_gb"`
	Status   types.String `tfsdk:"status"`
}

func (d *aiVolumesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ai_volumes"
}

func (d *aiVolumesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the storage volumes of a HiTechCloud AI Factory service " +
			"(`GET /api/service/{service_id}/volumes`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (same as `service_id`).",
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the AI Factory service.",
			},
			"volumes": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Volumes of the service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.StringAttribute{Computed: true, MarkdownDescription: "Volume identifier."},
						"name":       schema.StringAttribute{Computed: true, MarkdownDescription: "Volume name."},
						"cloud":      schema.StringAttribute{Computed: true, MarkdownDescription: "Cloud provider."},
						"region":     schema.StringAttribute{Computed: true, MarkdownDescription: "Region."},
						"size_in_gb": schema.Int64Attribute{Computed: true, MarkdownDescription: "Size in GB."},
						"status":     schema.StringAttribute{Computed: true, MarkdownDescription: "Volume status."},
					},
				},
			},
		},
	}
}

func (d *aiVolumesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *aiVolumesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data aiVolumesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := data.ServiceID.ValueString()
	volumes, err := d.client.ListAIVolumes(ctx, serviceID)
	if err != nil {
		resp.Diagnostics.AddError("Error listing volumes", err.Error())
		return
	}

	out := make([]aiVolumeItem, 0, len(volumes))
	for _, v := range volumes {
		out = append(out, aiVolumeItem{
			ID:       common.StringOrNull(v.ID),
			Name:     common.StringOrNull(v.Name),
			Cloud:    common.StringOrNull(v.Cloud),
			Region:   common.StringOrNull(v.Region),
			SizeInGB: common.Int64OrNull(v.SizeInGB),
			Status:   common.StringOrNull(v.Status),
		})
	}

	data.ID = types.StringValue(serviceID)
	data.Volumes = out
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
