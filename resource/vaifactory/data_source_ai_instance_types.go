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
	_ datasource.DataSource              = (*aiInstanceTypesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*aiInstanceTypesDataSource)(nil)
)

// AIInstanceTypesDataSource returns the hitechcloud_ai_instance_types data source.
func AIInstanceTypesDataSource() datasource.DataSource {
	return &aiInstanceTypesDataSource{}
}

type aiInstanceTypesDataSource struct {
	client *client.Client
}

type aiInstanceTypesModel struct {
	ID        types.String         `tfsdk:"id"`
	ServiceID types.String         `tfsdk:"service_id"`
	Types     []aiInstanceTypeItem `tfsdk:"types"`
}

type aiInstanceTypeItem struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	GPU   types.String `tfsdk:"gpu"`
	VCPUs types.Int64  `tfsdk:"vcpus"`
	RAM   types.String `tfsdk:"ram"`
}

func (d *aiInstanceTypesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ai_instance_types"
}

func (d *aiInstanceTypesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the GPU instance types available on a HiTechCloud AI " +
			"Factory service (`GET /api/service/{service_id}/instances/types`).",
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
				MarkdownDescription: "Available GPU instance types.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":    schema.StringAttribute{Computed: true, MarkdownDescription: "Instance type identifier."},
						"name":  schema.StringAttribute{Computed: true, MarkdownDescription: "Display name."},
						"gpu":   schema.StringAttribute{Computed: true, MarkdownDescription: "GPU model."},
						"vcpus": schema.Int64Attribute{Computed: true, MarkdownDescription: "Number of vCPUs."},
						"ram":   schema.StringAttribute{Computed: true, MarkdownDescription: "RAM specification."},
					},
				},
			},
		},
	}
}

func (d *aiInstanceTypesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *aiInstanceTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data aiInstanceTypesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := data.ServiceID.ValueString()
	items, err := d.client.ListAIInstanceTypes(ctx, serviceID)
	if err != nil {
		resp.Diagnostics.AddError("Error listing AI instance types", err.Error())
		return
	}

	out := make([]aiInstanceTypeItem, 0, len(items))
	for _, t := range items {
		out = append(out, aiInstanceTypeItem{
			ID:    common.StringOrNull(t.ID),
			Name:  common.StringOrNull(t.Name),
			GPU:   common.StringOrNull(t.GPU),
			VCPUs: common.Int64OrNull(t.VCPUs),
			RAM:   common.StringOrNull(t.RAM),
		})
	}

	data.ID = types.StringValue(serviceID)
	data.Types = out
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
