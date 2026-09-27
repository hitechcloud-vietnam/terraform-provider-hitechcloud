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
	_ datasource.DataSource              = (*aiClusterTypesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*aiClusterTypesDataSource)(nil)
)

// AIClusterTypesDataSource returns the hitechcloud_ai_cluster_types data source.
func AIClusterTypesDataSource() datasource.DataSource {
	return &aiClusterTypesDataSource{}
}

type aiClusterTypesDataSource struct {
	client *client.Client
}

type aiClusterTypesModel struct {
	ID        types.String        `tfsdk:"id"`
	ServiceID types.String        `tfsdk:"service_id"`
	Types     []aiClusterTypeItem `tfsdk:"types"`
}

type aiClusterTypeItem struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *aiClusterTypesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ai_cluster_types"
}

func (d *aiClusterTypesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the GPU cluster types available on a HiTechCloud AI " +
			"Factory service (`GET /api/service/{service_id}/clusters/types`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (same as `service_id`).",
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the AI Factory service.",
			},
			"types": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Available GPU cluster types.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true, MarkdownDescription: "Cluster type identifier."},
						"name": schema.StringAttribute{Computed: true, MarkdownDescription: "Display name."},
					},
				},
			},
		},
	}
}

func (d *aiClusterTypesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *aiClusterTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data aiClusterTypesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := data.ServiceID.ValueString()
	items, err := d.client.ListAIClusterTypes(ctx, serviceID)
	if err != nil {
		resp.Diagnostics.AddError("Error listing cluster types", err.Error())
		return
	}

	out := make([]aiClusterTypeItem, 0, len(items))
	for _, t := range items {
		out = append(out, aiClusterTypeItem{
			ID:   common.StringOrNull(t.ID),
			Name: common.StringOrNull(t.Name),
		})
	}

	data.ID = types.StringValue(serviceID)
	data.Types = out
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
