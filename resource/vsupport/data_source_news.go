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
	_ datasource.DataSource              = (*newsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*newsDataSource)(nil)
)

// NewsDataSource returns the hitechcloud_news data source.
func NewsDataSource() datasource.DataSource {
	return &newsDataSource{}
}

type newsDataSource struct {
	client *client.Client
}

type newsModel struct {
	ID   types.String `tfsdk:"id"`
	News []newsRow    `tfsdk:"news"`
	KB   []kbRow      `tfsdk:"knowledgebase"`
}

type newsRow struct {
	ID        types.String `tfsdk:"id"`
	Title     types.String `tfsdk:"title"`
	Body      types.String `tfsdk:"body"`
	CreatedAt types.String `tfsdk:"created_at"`
}

type kbRow struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Count types.Int64  `tfsdk:"article_count"`
}

func (d *newsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_news"
}

func (d *newsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads portal news and knowledgebase categories " +
			"(`GET /api/news`, `GET /api/knowledgebase`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Static identifier for this query (`news`).",
			},
			"news": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Latest news items.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.StringAttribute{Computed: true, MarkdownDescription: "News item ID."},
						"title":      schema.StringAttribute{Computed: true, MarkdownDescription: "News title."},
						"body":       schema.StringAttribute{Computed: true, MarkdownDescription: "News body text."},
						"created_at": schema.StringAttribute{Computed: true, MarkdownDescription: "Publication timestamp."},
					},
				},
			},
			"knowledgebase": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Knowledgebase categories.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":            schema.StringAttribute{Computed: true, MarkdownDescription: "Category ID."},
						"name":          schema.StringAttribute{Computed: true, MarkdownDescription: "Category name."},
						"article_count": schema.Int64Attribute{Computed: true, MarkdownDescription: "Number of articles in the category."},
					},
				},
			},
		},
	}
}

func (d *newsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *newsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data newsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID = types.StringValue("news")

	if items, err := d.client.ListNews(ctx); err == nil {
		data.News = make([]newsRow, 0, len(items))
		for _, item := range items {
			m, ok := client.AsMap(item)
			if !ok {
				continue
			}
			data.News = append(data.News, newsRow{
				ID:        common.StringOrNull(client.FirstString(m, "id", "news_id")),
				Title:     common.StringOrNull(client.FirstString(m, "title", "subject", "name")),
				Body:      common.StringOrNull(client.FirstString(m, "body", "content", "text")),
				CreatedAt: common.StringOrNull(client.FirstString(m, "created_at", "created", "date")),
			})
		}
	}

	if cats, err := d.client.ListKnowledgebaseCategories(ctx); err == nil {
		data.KB = make([]kbRow, 0, len(cats))
		for _, item := range cats {
			m, ok := client.AsMap(item)
			if !ok {
				continue
			}
			data.KB = append(data.KB, kbRow{
				ID:    common.StringOrNull(client.FirstString(m, "id", "category_id")),
				Name:  common.StringOrNull(client.FirstString(m, "name", "title", "label")),
				Count: types.Int64Value(client.FirstInt64(m, "article_count", "count", "articles")),
			})
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
