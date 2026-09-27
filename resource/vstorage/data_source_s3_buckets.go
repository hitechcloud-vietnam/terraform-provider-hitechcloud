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
	_ datasource.DataSource              = (*s3BucketsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*s3BucketsDataSource)(nil)
)

// S3BucketsDataSource returns the hitechcloud_s3_buckets data source.
func S3BucketsDataSource() datasource.DataSource {
	return &s3BucketsDataSource{}
}

type s3BucketsDataSource struct {
	client *client.Client
}

type s3BucketsModel struct {
	ID        types.String   `tfsdk:"id"`
	ServiceID types.String   `tfsdk:"service_id"`
	Buckets   []s3BucketItem `tfsdk:"buckets"`
}

type s3BucketItem struct {
	Name      types.String `tfsdk:"name"`
	Created   types.String `tfsdk:"created"`
	SizeBytes types.Int64  `tfsdk:"size_bytes"`
	Objects   types.Int64  `tfsdk:"objects"`
}

func (d *s3BucketsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_s3_buckets"
}

func (d *s3BucketsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the buckets of a HiTechCloud Ceph S3 service " +
			"(`GET /api/service/{service_id}/s3/buckets`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (same as `service_id`).",
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the S3 service.",
			},
			"buckets": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Buckets of the service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":       schema.StringAttribute{Computed: true, MarkdownDescription: "Bucket name."},
						"created":    schema.StringAttribute{Computed: true, MarkdownDescription: "Creation timestamp."},
						"size_bytes": schema.Int64Attribute{Computed: true, MarkdownDescription: "Total size in bytes."},
						"objects":    schema.Int64Attribute{Computed: true, MarkdownDescription: "Number of objects."},
					},
				},
			},
		},
	}
}

func (d *s3BucketsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *s3BucketsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data s3BucketsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := data.ServiceID.ValueString()
	buckets, err := d.client.ListS3Buckets(ctx, serviceID)
	if err != nil {
		resp.Diagnostics.AddError("Error listing S3 buckets", err.Error())
		return
	}

	out := make([]s3BucketItem, 0, len(buckets))
	for _, b := range buckets {
		out = append(out, s3BucketItem{
			Name:      common.StringOrNull(b.Name),
			Created:   common.StringOrNull(b.Created),
			SizeBytes: common.Int64OrNull(b.SizeBytes),
			Objects:   common.Int64OrNull(b.Objects),
		})
	}

	data.ID = types.StringValue(serviceID)
	data.Buckets = out
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
