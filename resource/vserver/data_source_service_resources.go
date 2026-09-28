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
	_ datasource.DataSource              = (*serviceResourcesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*serviceResourcesDataSource)(nil)
)

// ServiceResourcesDataSource returns the hitechcloud_service_resources data source.
func ServiceResourcesDataSource() datasource.DataSource {
	return &serviceResourcesDataSource{}
}

type serviceResourcesDataSource struct {
	client *client.Client
}

type serviceResourcesModel struct {
	ID             types.String        `tfsdk:"id"`
	ServiceID      types.String        `tfsdk:"service_id"`
	Resources      types.Map           `tfsdk:"resources"`
	UpgradeOptions []serviceUpgradeRow `tfsdk:"upgrade_options"`
	Networks       []serviceNetworkRow `tfsdk:"networks"`
	Images         []serviceImageRow   `tfsdk:"images"`
}

type serviceUpgradeRow struct {
	Name  types.String `tfsdk:"name"`
	Price types.String `tfsdk:"price"`
	Extra types.Map    `tfsdk:"extra"`
}

type serviceNetworkRow struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

type serviceImageRow struct {
	ID    types.String `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
	OS    types.String `tfsdk:"os"`
}

func (d *serviceResourcesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_resources"
}

func (d *serviceResourcesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads the resource envelope of a service: allocated resources, " +
			"available upgrade options, networks and images " +
			"(`GET /api/service/{id}/resources`, `/upgrade`, `/networks`, `/images`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The service ID (same as `service_id`).",
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HiTechCloud service ID.",
			},
			"resources": schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Resource allocations of the service."},
			"upgrade_options": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Upgrade targets available for the service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":  schema.StringAttribute{Computed: true, MarkdownDescription: "Upgrade target name."},
						"price": schema.StringAttribute{Computed: true, MarkdownDescription: "Upgrade price."},
						"extra": schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Remaining fields of the option."},
					},
				},
			},
			"networks": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Networks attached to the service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true, MarkdownDescription: "Network identifier."},
						"name": schema.StringAttribute{Computed: true, MarkdownDescription: "Network name."},
						"type": schema.StringAttribute{Computed: true, MarkdownDescription: "Network type."},
					},
				},
			},
			"images": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Images available to the service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":    schema.StringAttribute{Computed: true, MarkdownDescription: "Image identifier."},
						"label": schema.StringAttribute{Computed: true, MarkdownDescription: "Image label."},
						"os":    schema.StringAttribute{Computed: true, MarkdownDescription: "Operating system."},
					},
				},
			},
		},
	}
}

func (d *serviceResourcesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *serviceResourcesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data serviceResourcesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := data.ServiceID.ValueString()
	data.ID = types.StringValue(serviceID)

	if res, err := d.client.GetServiceResources(ctx, serviceID); err == nil {
		data.Resources = common.MapStrings(client.StringMap(res))
	}
	if opts, err := d.client.ListServiceUpgradeOptions(ctx, serviceID); err == nil {
		data.UpgradeOptions = make([]serviceUpgradeRow, 0, len(opts))
		for _, item := range opts {
			m, ok := client.AsMap(item)
			if !ok {
				continue
			}
			data.UpgradeOptions = append(data.UpgradeOptions, serviceUpgradeRow{
				Name:  types.StringValue(client.FirstString(m, "name", "title", "label")),
				Price: common.StringOrNull(client.FirstString(m, "price", "cost", "amount")),
				Extra: common.MapStrings(client.StringMap(m)),
			})
		}
	}
	if nets, err := d.client.ListServiceNetworks(ctx, serviceID); err == nil {
		data.Networks = make([]serviceNetworkRow, 0, len(nets))
		for _, item := range nets {
			m, ok := client.AsMap(item)
			if !ok {
				continue
			}
			data.Networks = append(data.Networks, serviceNetworkRow{
				ID:   common.StringOrNull(client.FirstString(m, "id", "network_id")),
				Name: common.StringOrNull(client.FirstString(m, "name", "label")),
				Type: common.StringOrNull(client.FirstString(m, "type", "kind")),
			})
		}
	}
	if imgs, err := d.client.ListServiceImages(ctx, serviceID); err == nil {
		data.Images = make([]serviceImageRow, 0, len(imgs))
		for _, item := range imgs {
			m, ok := client.AsMap(item)
			if !ok {
				continue
			}
			data.Images = append(data.Images, serviceImageRow{
				ID:    common.StringOrNull(client.FirstString(m, "id", "image", "file")),
				Label: common.StringOrNull(client.FirstString(m, "label", "name", "title")),
				OS:    common.StringOrNull(client.FirstString(m, "os", "distro", "version")),
			})
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
