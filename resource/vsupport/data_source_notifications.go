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
	_ datasource.DataSource              = (*notificationsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*notificationsDataSource)(nil)
)

// NotificationsDataSource returns the hitechcloud_notifications data source.
func NotificationsDataSource() datasource.DataSource {
	return &notificationsDataSource{}
}

type notificationsDataSource struct {
	client *client.Client
}

type notificationsModel struct {
	ID            types.String       `tfsdk:"id"`
	RelType       types.String       `tfsdk:"rel_type"`
	RelID         types.String       `tfsdk:"rel_id"`
	Notifications []notificationItem `tfsdk:"notifications"`
}

type notificationItem struct {
	ID      types.String `tfsdk:"id"`
	Title   types.String `tfsdk:"title"`
	Message types.String `tfsdk:"message"`
	Date    types.String `tfsdk:"date"`
	Status  types.String `tfsdk:"status"`
}

func (d *notificationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notifications"
}

func (d *notificationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the portal notifications of the HiTechCloud account " +
			"(`GET /api/notifications`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (`notifications`).",
			},
			"rel_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter by related object type.",
			},
			"rel_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Filter by related object identifier.",
			},
			"notifications": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Portal notifications.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":      schema.StringAttribute{Computed: true, MarkdownDescription: "Notification identifier."},
						"title":   schema.StringAttribute{Computed: true, MarkdownDescription: "Notification title."},
						"message": schema.StringAttribute{Computed: true, MarkdownDescription: "Notification message."},
						"date":    schema.StringAttribute{Computed: true, MarkdownDescription: "Creation timestamp."},
						"status":  schema.StringAttribute{Computed: true, MarkdownDescription: "Notification status."},
					},
				},
			},
		},
	}
}

func (d *notificationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *notificationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data notificationsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	relType := ""
	relID := ""
	if !data.RelType.IsNull() && !data.RelType.IsUnknown() {
		relType = data.RelType.ValueString()
	}
	if !data.RelID.IsNull() && !data.RelID.IsUnknown() {
		relID = data.RelID.ValueString()
	}

	notes, err := d.client.ListNotifications(ctx, relType, relID)
	if err != nil {
		resp.Diagnostics.AddError("Error listing notifications", err.Error())
		return
	}

	items := make([]notificationItem, 0, len(notes))
	for _, n := range notes {
		items = append(items, notificationItem{
			ID:      common.StringOrNull(n.ID),
			Title:   common.StringOrNull(n.Title),
			Message: common.StringOrNull(n.Message),
			Date:    common.StringOrNull(n.Date),
			Status:  common.StringOrNull(n.Status),
		})
	}

	data.ID = types.StringValue("notifications")
	data.Notifications = items
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
