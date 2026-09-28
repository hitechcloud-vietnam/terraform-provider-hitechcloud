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
	_ datasource.DataSource              = (*s3ConnectionDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*s3ConnectionDataSource)(nil)
)

// S3ConnectionDataSource returns the hitechcloud_s3_connection data source.
func S3ConnectionDataSource() datasource.DataSource {
	return &s3ConnectionDataSource{}
}

type s3ConnectionDataSource struct {
	client *client.Client
}

type s3ConnectionModel struct {
	ID        types.String `tfsdk:"id"`
	ServiceID types.String `tfsdk:"service_id"`
	Endpoint  types.String `tfsdk:"endpoint"`
	Region    types.String `tfsdk:"region"`
	AccessKey types.String `tfsdk:"access_key"`
	SecretKey types.String `tfsdk:"secret_key"`
	Usage     types.Map    `tfsdk:"usage"`
	Metrics   types.Map    `tfsdk:"metrics"`
}

func (d *s3ConnectionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_connection"
}

func (d *s3ConnectionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads an S3 object-storage service: connection info, credentials " +
			"and usage metrics (`GET /api/service/{id}/s3`, `/s3/credentials`, `/s3/usage`, `/s3/metrics`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The service ID (same as `service_id`).",
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HiTechCloud service ID (`hb_accounts.id`) of the S3 service.",
			},
			"endpoint":   schema.StringAttribute{Computed: true, MarkdownDescription: "S3 endpoint URL."},
			"region":     schema.StringAttribute{Computed: true, MarkdownDescription: "S3 region."},
			"access_key": schema.StringAttribute{Computed: true, Sensitive: true, MarkdownDescription: "S3 access key."},
			"secret_key": schema.StringAttribute{Computed: true, Sensitive: true, MarkdownDescription: "S3 secret key."},
			"usage":      schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Billing-relevant usage counters with units."},
			"metrics":    schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Raw metric values keyed by metric name."},
		},
	}
}

func (d *s3ConnectionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *s3ConnectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data s3ConnectionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := data.ServiceID.ValueString()
	data.ID = types.StringValue(serviceID)

	if info, err := d.client.GetS3ConnectionInfo(ctx, serviceID); err == nil {
		data.Endpoint = common.StringOrNull(client.FirstString(info, "endpoint", "url", "host"))
		data.Region = common.StringOrNull(client.FirstString(info, "region"))
	}
	if creds, err := d.client.GetS3Credentials(ctx, serviceID); err == nil {
		data.AccessKey = common.StringOrNull(client.FirstString(creds, "access_key", "accesskey", "key"))
		data.SecretKey = common.StringOrNull(client.FirstString(creds, "secret_key", "secretkey", "secret"))
		if data.Endpoint.IsNull() {
			data.Endpoint = common.StringOrNull(client.FirstString(creds, "endpoint", "url", "host"))
		}
	}
	if usage, err := d.client.GetS3Usage(ctx, serviceID); err == nil {
		data.Usage = common.MapStrings(client.StringMap(usage))
	}
	if metrics, err := d.client.GetS3Metrics(ctx, serviceID, ""); err == nil {
		data.Metrics = common.MapStrings(client.StringMap(metrics))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
