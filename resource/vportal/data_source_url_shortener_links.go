// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vportal

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
	_ datasource.DataSource              = (*urlShortenerLinksDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*urlShortenerLinksDataSource)(nil)
)

// URLShortenerLinksDataSource returns the hitechcloud_url_shortener_links data source.
func URLShortenerLinksDataSource() datasource.DataSource {
	return &urlShortenerLinksDataSource{}
}

type urlShortenerLinksDataSource struct {
	client *client.Client
}

type urlShortenerLinksModel struct {
	ID      types.String  `tfsdk:"id"`
	Page    types.Int64   `tfsdk:"page"`
	PerPage types.Int64   `tfsdk:"per_page"`
	Search  types.String  `tfsdk:"search"`
	Links   []urlLinkItem `tfsdk:"links"`
}

type urlLinkItem struct {
	ID       types.String `tfsdk:"id"`
	URL      types.String `tfsdk:"url"`
	ShortURL types.String `tfsdk:"short_url"`
	Label    types.String `tfsdk:"label"`
	Clicks   types.Int64  `tfsdk:"clicks"`
}

func (d *urlShortenerLinksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_url_shortener_links"
}

func (d *urlShortenerLinksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the shortened URLs of the HiTechCloud portal " +
			"(`GET /api/url-shortener/links`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (`url_shortener_links`).",
			},
			"page": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Page number. Defaults to `1`.",
			},
			"per_page": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Results per page.",
			},
			"search": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Search term.",
			},
			"links": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Shortened links.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":        schema.StringAttribute{Computed: true, MarkdownDescription: "Link identifier."},
						"url":       schema.StringAttribute{Computed: true, MarkdownDescription: "Target URL."},
						"short_url": schema.StringAttribute{Computed: true, MarkdownDescription: "Short URL."},
						"label":     schema.StringAttribute{Computed: true, MarkdownDescription: "Label / alias."},
						"clicks":    schema.Int64Attribute{Computed: true, MarkdownDescription: "Click counter."},
					},
				},
			},
		},
	}
}

func (d *urlShortenerLinksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *urlShortenerLinksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data urlShortenerLinksModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var page, perPage int64 = 1, 0
	if !data.Page.IsNull() && !data.Page.IsUnknown() {
		page = data.Page.ValueInt64()
	}
	if !data.PerPage.IsNull() && !data.PerPage.IsUnknown() {
		perPage = data.PerPage.ValueInt64()
	}
	search := ""
	if !data.Search.IsNull() && !data.Search.IsUnknown() {
		search = data.Search.ValueString()
	}

	links, err := d.client.ListURLLinks(ctx, page, perPage, search)
	if err != nil {
		resp.Diagnostics.AddError("Error listing short links", err.Error())
		return
	}

	items := make([]urlLinkItem, 0, len(links))
	for _, l := range links {
		items = append(items, urlLinkItem{
			ID:       common.StringOrNull(l.ID),
			URL:      common.StringOrNull(l.URL),
			ShortURL: common.StringOrNull(l.Short),
			Label:    common.StringOrNull(l.Label),
			Clicks:   types.Int64Value(l.Clicks),
		})
	}

	data.ID = types.StringValue("url_shortener_links")
	data.Links = items
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
