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
	_ datasource.DataSource              = (*ipamDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*ipamDataSource)(nil)
)

// IPAMDataSource returns the hitechcloud_ipam data source.
func IPAMDataSource() datasource.DataSource {
	return &ipamDataSource{}
}

type ipamDataSource struct {
	client *client.Client
}

type ipamModel struct {
	ID        types.String   `tfsdk:"id"`
	ServiceID types.String   `tfsdk:"service_id"`
	IPs       []ipamIP       `tfsdk:"ips"`
	Subnets   []ipamSubnet   `tfsdk:"subnets"`
	RDNS      []ipamRDNSItem `tfsdk:"rdns"`
}

type ipamIP struct {
	IP      types.String `tfsdk:"ip"`
	Subnet  types.String `tfsdk:"subnet"`
	Status  types.String `tfsdk:"status"`
	Comment types.String `tfsdk:"comment"`
}

type ipamSubnet struct {
	CIDR     types.String `tfsdk:"cidr"`
	Gateway  types.String `tfsdk:"gateway"`
	Type     types.String `tfsdk:"type"`
	Capacity types.Int64  `tfsdk:"capacity"`
}

type ipamRDNSItem struct {
	IP       types.String `tfsdk:"ip"`
	Hostname types.String `tfsdk:"hostname"`
}

func (d *ipamDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipam"
}

func (d *ipamDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads an IPAM (IP address management) service: allocated IPs, " +
			"subnets and reverse DNS records " +
			"(`GET /api/service/{id}/htcipam/ips`, `/htcipam/subnets`, `/htcipam/rdns`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The service ID (same as `service_id`).",
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HiTechCloud service ID of the IPAM service.",
			},
			"ips": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "IP addresses allocated from the IPAM pool.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip":      schema.StringAttribute{Computed: true, MarkdownDescription: "IP address."},
						"subnet":  schema.StringAttribute{Computed: true, MarkdownDescription: "Parent subnet."},
						"status":  schema.StringAttribute{Computed: true, MarkdownDescription: "Allocation status."},
						"comment": schema.StringAttribute{Computed: true, MarkdownDescription: "Operator comment."},
					},
				},
			},
			"subnets": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Subnets managed by the IPAM service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"cidr":     schema.StringAttribute{Computed: true, MarkdownDescription: "Subnet in CIDR notation."},
						"gateway":  schema.StringAttribute{Computed: true, MarkdownDescription: "Gateway address."},
						"type":     schema.StringAttribute{Computed: true, MarkdownDescription: "Address family (ipv4/ipv6)."},
						"capacity": schema.Int64Attribute{Computed: true, MarkdownDescription: "Number of addresses in the subnet."},
					},
				},
			},
			"rdns": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Reverse DNS records of the service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ip":       schema.StringAttribute{Computed: true, MarkdownDescription: "IP address."},
						"hostname": schema.StringAttribute{Computed: true, MarkdownDescription: "PTR hostname."},
					},
				},
			},
		},
	}
}

func (d *ipamDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ipamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ipamModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := data.ServiceID.ValueString()
	data.ID = types.StringValue(serviceID)

	if ips, err := d.client.ListIPAMIPs(ctx, serviceID); err == nil {
		data.IPs = make([]ipamIP, 0, len(ips))
		for _, item := range ips {
			m, ok := client.AsMap(item)
			if !ok {
				continue
			}
			data.IPs = append(data.IPs, ipamIP{
				IP:      types.StringValue(client.FirstString(m, "ip", "address", "addr")),
				Subnet:  common.StringOrNull(client.FirstString(m, "subnet", "cidr", "network")),
				Status:  common.StringOrNull(client.FirstString(m, "status", "state")),
				Comment: common.StringOrNull(client.FirstString(m, "comment", "note", "description")),
			})
		}
	}

	if subnets, err := d.client.ListIPAMSubnets(ctx, serviceID); err == nil {
		data.Subnets = make([]ipamSubnet, 0, len(subnets))
		for _, item := range subnets {
			m, ok := client.AsMap(item)
			if !ok {
				continue
			}
			data.Subnets = append(data.Subnets, ipamSubnet{
				CIDR:     types.StringValue(client.FirstString(m, "cidr", "subnet", "network")),
				Gateway:  common.StringOrNull(client.FirstString(m, "gateway", "gw")),
				Type:     common.StringOrNull(client.FirstString(m, "type", "family", "version")),
				Capacity: types.Int64Value(client.FirstInt64(m, "capacity", "size", "count")),
			})
		}
	}

	records, err := d.client.GetIPAMRDNS(ctx, serviceID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading IPAM reverse DNS", err.Error())
		return
	}
	data.RDNS = make([]ipamRDNSItem, 0, len(records))
	for _, r := range records {
		data.RDNS = append(data.RDNS, ipamRDNSItem{
			IP:       types.StringValue(r.IP),
			Hostname: common.StringOrNull(r.Hostname),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
