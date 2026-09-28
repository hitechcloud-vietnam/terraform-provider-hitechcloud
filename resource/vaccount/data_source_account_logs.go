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
	_ datasource.DataSource              = (*accountLogsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*accountLogsDataSource)(nil)
)

// AccountLogsDataSource returns the hitechcloud_account_logs data source.
func AccountLogsDataSource() datasource.DataSource {
	return &accountLogsDataSource{}
}

type accountLogsDataSource struct {
	client *client.Client
}

type accountLogsModel struct {
	ID   types.String `tfsdk:"id"`
	Logs []accountLog `tfsdk:"logs"`
}

type accountLog struct {
	ID        types.String `tfsdk:"id"`
	Action    types.String `tfsdk:"action"`
	IP        types.String `tfsdk:"ip"`
	UserAgent types.String `tfsdk:"user_agent"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func (d *accountLogsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_logs"
}

func (d *accountLogsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads the account activity log of the authenticated user " +
			"(`GET /api/logs`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Static identifier for this query.",
			},
			"logs": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Account activity log entries, newest first.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.StringAttribute{Computed: true, MarkdownDescription: "Log entry ID."},
						"action":     schema.StringAttribute{Computed: true, MarkdownDescription: "Action performed."},
						"ip":         schema.StringAttribute{Computed: true, MarkdownDescription: "Source IP address."},
						"user_agent": schema.StringAttribute{Computed: true, MarkdownDescription: "Client user agent."},
						"created_at": schema.StringAttribute{Computed: true, MarkdownDescription: "Timestamp of the entry."},
					},
				},
			},
		},
	}
}

func (d *accountLogsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *accountLogsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data accountLogsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	logs, err := d.client.GetAccountLogs(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading account logs", err.Error())
		return
	}
	data.ID = types.StringValue("account_logs")
	data.Logs = make([]accountLog, 0, len(logs))
	for _, item := range logs {
		m, ok := client.AsMap(item)
		if !ok {
			continue
		}
		data.Logs = append(data.Logs, accountLog{
			ID:        common.StringOrNull(client.FirstString(m, "id", "log_id")),
			Action:    common.StringOrNull(client.FirstString(m, "action", "event", "type")),
			IP:        common.StringOrNull(client.FirstString(m, "ip", "ip_address", "ipaddr")),
			UserAgent: common.StringOrNull(client.FirstString(m, "user_agent", "useragent", "ua")),
			CreatedAt: common.StringOrNull(client.FirstString(m, "created_at", "created", "time", "date")),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
