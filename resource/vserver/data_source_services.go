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
	_ datasource.DataSource              = (*servicesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*servicesDataSource)(nil)
)

// ServicesDataSource returns the hitechcloud_services data source.
func ServicesDataSource() datasource.DataSource {
	return &servicesDataSource{}
}

type servicesDataSource struct {
	client *client.Client
}

type servicesModel struct {
	ID       types.String  `tfsdk:"id"`
	Services []serviceItem `tfsdk:"services"`
}

// serviceItem is the shared nested object for services.
type serviceItem struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Status       types.String `tfsdk:"status"`
	Label        types.String `tfsdk:"label"`
	Domain       types.String `tfsdk:"domain"`
	Group        types.String `tfsdk:"group"`
	Cycle        types.String `tfsdk:"cycle"`
	Amount       types.String `tfsdk:"amount"`
	NextDue      types.String `tfsdk:"next_due"`
	IP           types.String `tfsdk:"ip"`
	RegisterDate types.String `tfsdk:"register_date"`
}

// serviceItemSchema mirrors the nested object for the single-service data source.
func serviceItemAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":            schema.StringAttribute{Computed: true, MarkdownDescription: "Service identifier (HostBill service id)."},
		"name":          schema.StringAttribute{Computed: true, MarkdownDescription: "Product / service name."},
		"status":        schema.StringAttribute{Computed: true, MarkdownDescription: "Service status."},
		"label":         schema.StringAttribute{Computed: true, MarkdownDescription: "Custom label."},
		"domain":        schema.StringAttribute{Computed: true, MarkdownDescription: "Primary domain / hostname."},
		"group":         schema.StringAttribute{Computed: true, MarkdownDescription: "Product group."},
		"cycle":         schema.StringAttribute{Computed: true, MarkdownDescription: "Billing cycle."},
		"amount":        schema.StringAttribute{Computed: true, MarkdownDescription: "Recurring amount."},
		"next_due":      schema.StringAttribute{Computed: true, MarkdownDescription: "Next due date."},
		"ip":            schema.StringAttribute{Computed: true, MarkdownDescription: "Dedicated IP."},
		"register_date": schema.StringAttribute{Computed: true, MarkdownDescription: "Registration date."},
	}
}

func toServiceItem(s client.Service) serviceItem {
	return serviceItem{
		ID:           common.StringOrNull(s.ID),
		Name:         common.StringOrNull(s.Name),
		Status:       common.StringOrNull(s.Status),
		Label:        common.StringOrNull(s.Label),
		Domain:       common.StringOrNull(s.Domain),
		Group:        common.StringOrNull(s.Group),
		Cycle:        common.StringOrNull(s.Cycle),
		Amount:       common.StringOrNull(s.Amount),
		NextDue:      common.StringOrNull(s.NextDue),
		IP:           common.StringOrNull(s.IP),
		RegisterDate: common.StringOrNull(s.RegisterDate),
	}
}

func (d *servicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_services"
}

func (d *servicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists all services of the HiTechCloud account (`GET /api/service`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (`services`).",
			},
			"services": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Services of the account.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: serviceItemAttributes(),
				},
			},
		},
	}
}

func (d *servicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *servicesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	services, err := d.client.ListServices(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error listing services", err.Error())
		return
	}

	items := make([]serviceItem, 0, len(services))
	for _, s := range services {
		items = append(items, toServiceItem(s))
	}

	data := servicesModel{
		ID:       types.StringValue("services"),
		Services: items,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
