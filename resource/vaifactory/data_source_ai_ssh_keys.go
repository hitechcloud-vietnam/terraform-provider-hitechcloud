// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vaifactory

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
	_ datasource.DataSource              = (*aiSSHKeysDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*aiSSHKeysDataSource)(nil)
)

// AISSHKeysDataSource returns the hitechcloud_ai_ssh_keys data source.
func AISSHKeysDataSource() datasource.DataSource {
	return &aiSSHKeysDataSource{}
}

type aiSSHKeysDataSource struct {
	client *client.Client
}

type aiSSHKeysModel struct {
	ID        types.String   `tfsdk:"id"`
	ServiceID types.String   `tfsdk:"service_id"`
	Keys      []aiSSHKeyItem `tfsdk:"keys"`
}

type aiSSHKeyItem struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	PublicKey types.String `tfsdk:"public_key"`
	IsDefault types.Bool   `tfsdk:"is_default"`
}

func (d *aiSSHKeysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ai_ssh_keys"
}

func (d *aiSSHKeysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists the SSH keys registered with a HiTechCloud AI Factory " +
			"service (`GET /api/service/{service_id}/sshkeys`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Data source identifier (same as `service_id`).",
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the AI Factory service.",
			},
			"keys": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "SSH keys of the service.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.StringAttribute{Computed: true, MarkdownDescription: "Key identifier."},
						"name":       schema.StringAttribute{Computed: true, MarkdownDescription: "Key name."},
						"public_key": schema.StringAttribute{Computed: true, MarkdownDescription: "Public key material."},
						"is_default": schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether this is the default key."},
					},
				},
			},
		},
	}
}

func (d *aiSSHKeysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *aiSSHKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data aiSSHKeysModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID := data.ServiceID.ValueString()
	keys, err := d.client.ListAISSHKeys(ctx, serviceID)
	if err != nil {
		resp.Diagnostics.AddError("Error listing SSH keys", err.Error())
		return
	}

	out := make([]aiSSHKeyItem, 0, len(keys))
	for _, k := range keys {
		out = append(out, aiSSHKeyItem{
			ID:        common.StringOrNull(k.ID),
			Name:      common.StringOrNull(k.Name),
			PublicKey: common.StringOrNull(k.PublicKey),
			IsDefault: types.BoolValue(k.IsDefault),
		})
	}

	data.ID = types.StringValue(serviceID)
	data.Keys = out
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
