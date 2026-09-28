// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vaccount

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
	_ datasource.DataSource              = (*locationsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*locationsDataSource)(nil)
)

// LocationsDataSource returns the hitechcloud_locations data source.
func LocationsDataSource() datasource.DataSource {
	return &locationsDataSource{}
}

type locationsDataSource struct {
	client *client.Client
}

type locationsModel struct {
	ID        types.String    `tfsdk:"id"`
	Countries []locationEntry `tfsdk:"countries"`
}

type locationEntry struct {
	ID       types.String    `tfsdk:"id"`
	Name     types.String    `tfsdk:"name"`
	Code     types.String    `tfsdk:"code"`
	Children []locationChild `tfsdk:"children"`
}

type locationChild struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *locationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_locations"
}

func (d *locationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads the LocationV2 geography hierarchy: countries and their " +
			"states (`GET /api/location_v2/countries`, `/location_v2/{code}/state`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Static identifier for this query (`locations`).",
			},
			"countries": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Countries with their states.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.StringAttribute{Computed: true, MarkdownDescription: "Country identifier."},
						"name": schema.StringAttribute{Computed: true, MarkdownDescription: "Country name."},
						"code": schema.StringAttribute{Computed: true, MarkdownDescription: "Country code."},
						"children": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "States / provinces of the country.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id":   schema.StringAttribute{Computed: true, MarkdownDescription: "State identifier."},
									"name": schema.StringAttribute{Computed: true, MarkdownDescription: "State name."},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *locationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *locationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data locationsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	countries, err := d.client.ListCountries(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading countries", err.Error())
		return
	}
	data.ID = types.StringValue("locations")
	data.Countries = make([]locationEntry, 0, len(countries))
	for _, item := range countries {
		m, ok := client.AsMap(item)
		if !ok {
			continue
		}
		entry := locationEntry{
			ID:   common.StringOrNull(client.FirstString(m, "id", "code")),
			Name: common.StringOrNull(client.FirstString(m, "name", "country", "title")),
			Code: common.StringOrNull(client.FirstString(m, "code", "iso", "short")),
		}
		if code := client.FirstString(m, "code", "iso", "short", "id"); code != "" {
			if states, err := d.client.ListStates(ctx, code); err == nil {
				entry.Children = make([]locationChild, 0, len(states))
				for _, s := range states {
					sm, ok := client.AsMap(s)
					if !ok {
						continue
					}
					entry.Children = append(entry.Children, locationChild{
						ID:   common.StringOrNull(client.FirstString(sm, "id", "code")),
						Name: common.StringOrNull(client.FirstString(sm, "name", "state", "title")),
					})
				}
			}
		}
		data.Countries = append(data.Countries, entry)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
