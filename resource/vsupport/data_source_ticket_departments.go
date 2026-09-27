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
	_ datasource.DataSource              = (*ticketDepartmentsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*ticketDepartmentsDataSource)(nil)
)

// TicketDepartmentsDataSource returns the hitechcloud_ticket_departments data source.
func TicketDepartmentsDataSource() datasource.DataSource {
	return &ticketDepartmentsDataSource{}
}

type ticketDepartmentsDataSource struct {
	client *client.Client
}

type ticketDepartmentsModel struct {
	ID          types.String           `tfsdk:"id"`
	Departments []ticketDepartmentItem `tfsdk:"departments"`
}

type ticketDepartmentItem struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *ticketDepartmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ticket_departments"
}

func (d *ticketDepartmentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the support departments of the HiTechCloud account " +
			"(`GET /api/ticket/departments`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (`ticket_departments`).",
			},
			"departments": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Support departments.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true, MarkdownDescription: "Department identifier."},
						"name": schema.StringAttribute{Computed: true, MarkdownDescription: "Department name."},
					},
				},
			},
		},
	}
}

func (d *ticketDepartmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ticketDepartmentsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	depts, err := d.client.ListTicketDepartments(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing ticket departments", err.Error())
		return
	}

	items := make([]ticketDepartmentItem, 0, len(depts))
	for _, dept := range depts {
		items = append(items, ticketDepartmentItem{
			ID:   common.StringOrNull(dept.ID),
			Name: common.StringOrNull(dept.Name),
		})
	}

	data := ticketDepartmentsModel{
		ID:          types.StringValue("ticket_departments"),
		Departments: items,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
