// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vdns

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
	_ datasource.DataSource              = (*dnsZonesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*dnsZonesDataSource)(nil)
)

// DNSZonesDataSource returns the hitechcloud_dns_zones data source.
func DNSZonesDataSource() datasource.DataSource {
	return &dnsZonesDataSource{}
}

type dnsZonesDataSource struct {
	client *client.Client
}

type dnsZonesModel struct {
	ID        types.String    `tfsdk:"id"`
	ServiceID types.String    `tfsdk:"service_id"`
	Zones     []dnsZoneRecord `tfsdk:"zones"`
}

type dnsZoneRecord struct {
	ID      types.String    `tfsdk:"id"`
	Name    types.String    `tfsdk:"name"`
	Records []dnsRecordItem `tfsdk:"records"`
}

// dnsRecordItemSchema mirrors dnsRecordSchema for data source usage (the
// datasource/schema package has its own nested attribute types).
func dnsRecordItemSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Record identifier inside the zone.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Record name (host).",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Record type (A, AAAA, CNAME, MX, TXT, ...).",
			},
			"content": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Record value.",
			},
			"ttl": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Time to live in seconds.",
			},
			"priority": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Priority (MX/SRV records).",
			},
		},
	}
}

func (d *dnsZonesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_zones"
}

func (d *dnsZonesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the DNS zones of a HiTechCloud DNS service " +
			"(`GET /api/service/{service_id}/dns`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `<service_id>` (data source marker).",
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the DNS service.",
			},
			"zones": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "DNS zones of the service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Zone identifier.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Zone name (domain).",
						},
						"records": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "Records in the zone.",
							NestedObject:        dnsRecordItemSchema(),
						},
					},
				},
			},
		},
	}
}

func (d *dnsZonesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *dnsZonesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dnsZonesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := data.ServiceID.ValueString()
	zones, err := d.client.ListDNSZones(ctx, serviceID)
	if err != nil {
		resp.Diagnostics.AddError("Error listing DNS zones", err.Error())
		return
	}

	data.ID = types.StringValue(client.FormatID(serviceID))
	items := make([]dnsZoneRecord, 0, len(zones))
	for _, z := range zones {
		// Fetch the records listing for each zone (detailed read).
		if full, err := d.client.GetDNSZone(ctx, serviceID, z.ID); err == nil {
			z = *full
		}
		items = append(items, dnsZoneRecord{
			ID:      common.StringOrNull(z.ID),
			Name:    common.StringOrNull(z.Name),
			Records: toRecordItems(z.Records),
		})
	}
	data.Zones = items

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
