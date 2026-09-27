// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/client"
)

var (
	_ datasource.DataSource              = (*serviceDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*serviceDataSource)(nil)
)

// ServiceDataSource returns the hitechcloud_service data source.
func ServiceDataSource() datasource.DataSource {
	return &serviceDataSource{}
}

type serviceDataSource struct {
	client *client.Client
}

type serviceModel struct {
	ID           types.String `tfsdk:"id"`
	ServiceID    types.String `tfsdk:"service_id"`
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

func (d *serviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (d *serviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads a single HiTechCloud service (`GET /api/service/{service_id}`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (same as `service_id`).",
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id to look up.",
			},
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
		},
	}
}

func (d *serviceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *serviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data serviceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := data.ServiceID.ValueString()
	svc, err := d.client.GetService(ctx, serviceID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading service", err.Error())
		return
	}

	item := toServiceItem(*svc)
	data.ID = types.StringValue(serviceID)
	data.Name = item.Name
	data.Status = item.Status
	data.Label = item.Label
	data.Domain = item.Domain
	data.Group = item.Group
	data.Cycle = item.Cycle
	data.Amount = item.Amount
	data.NextDue = item.NextDue
	data.IP = item.IP
	data.RegisterDate = item.RegisterDate

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
