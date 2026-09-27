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
	_ datasource.DataSource              = (*ticketsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*ticketsDataSource)(nil)
)

// TicketsDataSource returns the hitechcloud_tickets data source.
func TicketsDataSource() datasource.DataSource {
	return &ticketsDataSource{}
}

type ticketsDataSource struct {
	client *client.Client
}

type ticketsModel struct {
	ID      types.String `tfsdk:"id"`
	Tickets []ticketItem `tfsdk:"tickets"`
}

type ticketItem struct {
	Number      types.String `tfsdk:"number"`
	Subject     types.String `tfsdk:"subject"`
	Status      types.String `tfsdk:"status"`
	Department  types.String `tfsdk:"department"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

func (d *ticketsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tickets"
}

func (d *ticketsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the support tickets of the HiTechCloud account " +
			"(`GET /api/tickets`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (`tickets`).",
			},
			"tickets": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Support tickets.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"number":       schema.StringAttribute{Computed: true, MarkdownDescription: "Ticket number."},
						"subject":      schema.StringAttribute{Computed: true, MarkdownDescription: "Ticket subject."},
						"status":       schema.StringAttribute{Computed: true, MarkdownDescription: "Ticket status."},
						"department":   schema.StringAttribute{Computed: true, MarkdownDescription: "Support department."},
						"last_updated": schema.StringAttribute{Computed: true, MarkdownDescription: "Last update timestamp."},
					},
				},
			},
		},
	}
}

func (d *ticketsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ticketsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tickets, err := d.client.ListTickets(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing tickets", err.Error())
		return
	}

	items := make([]ticketItem, 0, len(tickets))
	for _, t := range tickets {
		items = append(items, ticketItem{
			Number:      common.StringOrNull(t.Number),
			Subject:     common.StringOrNull(t.Subject),
			Status:      common.StringOrNull(t.Status),
			Department:  common.StringOrNull(t.Department),
			LastUpdated: common.StringOrNull(t.LastUpdated),
		})
	}

	data := ticketsModel{
		ID:      types.StringValue("tickets"),
		Tickets: items,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
