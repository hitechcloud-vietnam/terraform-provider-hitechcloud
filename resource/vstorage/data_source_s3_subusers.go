// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vstorage

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
	_ datasource.DataSource              = (*s3SubusersDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*s3SubusersDataSource)(nil)
)

// S3SubusersDataSource returns the hitechcloud_s3_subusers data source.
func S3SubusersDataSource() datasource.DataSource {
	return &s3SubusersDataSource{}
}

type s3SubusersDataSource struct {
	client *client.Client
}

type s3SubusersModel struct {
	ID        types.String    `tfsdk:"id"`
	ServiceID types.String    `tfsdk:"service_id"`
	Subusers  []s3SubuserItem `tfsdk:"subusers"`
}

type s3SubuserItem struct {
	Name      types.String `tfsdk:"name"`
	Access    types.String `tfsdk:"access"`
	AccessKey types.String `tfsdk:"access_key"`
}

func (d *s3SubusersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_subusers"
}

func (d *s3SubusersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the sub-users of a HiTechCloud Ceph S3 service " +
			"(`GET /api/service/{service_id}/s3/subusers`). Secret keys are never " +
			"returned by the API in listings.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (same as `service_id`).",
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the S3 service.",
			},
			"subusers": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Sub-users of the service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":       schema.StringAttribute{Computed: true, MarkdownDescription: "Sub-user name."},
						"access":     schema.StringAttribute{Computed: true, MarkdownDescription: "Permission level."},
						"access_key": schema.StringAttribute{Computed: true, MarkdownDescription: "Access key id."},
					},
				},
			},
		},
	}
}

func (d *s3SubusersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *s3SubusersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data s3SubusersModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := data.ServiceID.ValueString()
	users, err := d.client.ListS3Subusers(ctx, serviceID)
	if err != nil {
		resp.Diagnostics.AddError("Error listing S3 sub-users", err.Error())
		return
	}

	out := make([]s3SubuserItem, 0, len(users))
	for _, u := range users {
		out = append(out, s3SubuserItem{
			Name:      common.StringOrNull(u.Name),
			Access:    common.StringOrNull(u.Access),
			AccessKey: common.StringOrNull(u.AccessKey),
		})
	}

	data.ID = types.StringValue(serviceID)
	data.Subusers = out
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
