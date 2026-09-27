// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vsupport

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
	_ datasource.DataSource              = (*statusesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*statusesDataSource)(nil)
)

// StatusesDataSource returns the hitechcloud_statuses data source.
func StatusesDataSource() datasource.DataSource {
	return &statusesDataSource{}
}

type statusesDataSource struct {
	client *client.Client
}

type statusesModel struct {
	ID       types.String `tfsdk:"id"`
	Statuses []statusItem `tfsdk:"statuses"`
}

type statusItem struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Status types.String `tfsdk:"status"`
	Type   types.String `tfsdk:"type"`
}

func (d *statusesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_statuses"
}

func (d *statusesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the service status entries of the HiTechCloud account " +
			"(`GET /api/statuses`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (`statuses`).",
			},
			"statuses": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Service status entries.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":     schema.StringAttribute{Computed: true, MarkdownDescription: "Service identifier."},
						"name":   schema.StringAttribute{Computed: true, MarkdownDescription: "Service name."},
						"status": schema.StringAttribute{Computed: true, MarkdownDescription: "Current status."},
						"type":   schema.StringAttribute{Computed: true, MarkdownDescription: "Service type / group."},
					},
				},
			},
		},
	}
}

func (d *statusesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *statusesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	entries, err := d.client.ListStatuses(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing statuses", err.Error())
		return
	}

	items := make([]statusItem, 0, len(entries))
	for _, e := range entries {
		items = append(items, statusItem{
			ID:     common.StringOrNull(e.ID),
			Name:   common.StringOrNull(e.Name),
			Status: common.StringOrNull(e.Status),
			Type:   common.StringOrNull(e.Type),
		})
	}

	data := statusesModel{
		ID:       types.StringValue("statuses"),
		Statuses: items,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
