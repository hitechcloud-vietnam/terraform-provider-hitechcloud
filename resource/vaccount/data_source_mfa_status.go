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
	_ datasource.DataSource              = (*mfaStatusDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*mfaStatusDataSource)(nil)
)

// MFAStatusDataSource returns the hitechcloud_mfa_status data source.
func MFAStatusDataSource() datasource.DataSource {
	return &mfaStatusDataSource{}
}

type mfaStatusDataSource struct {
	client *client.Client
}

type mfaStatusModel struct {
	ID                 types.String `tfsdk:"id"`
	UserType           types.String `tfsdk:"user_type"`
	UserID             types.String `tfsdk:"user_id"`
	PasskeyStatus      types.Map    `tfsdk:"passkey_status"`
	EmailMFAStatus     types.Map    `tfsdk:"email_mfa_status"`
	PasskeyCredentials []passkeyRow `tfsdk:"passkey_credentials"`
	EmailMFACodes      []mfaCodeRow `tfsdk:"email_mfa_codes"`
}

type passkeyRow struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	CreatedAt  types.String `tfsdk:"created_at"`
	LastUsedAt types.String `tfsdk:"last_used_at"`
}

type mfaCodeRow struct {
	ID        types.String `tfsdk:"id"`
	Purpose   types.String `tfsdk:"purpose"`
	Status    types.String `tfsdk:"status"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func (d *mfaStatusDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mfa_status"
}

func (d *mfaStatusDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads MFA state of a user: passkey status and credentials, " +
			"email MFA status and codes " +
			"(`GET /api/passkeyv2/status/{user_type}/{user_id}`, " +
			"`/passkeyv2/credentials/{user_type}/{user_id}`, " +
			"`/email_mfa_v2/status/{user_type}/{user_id}`, " +
			"`/email_mfa_v2/list/{user_type}/{user_id}`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `user_type/user_id`.",
			},
			"user_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "User type, for example `client`.",
			},
			"user_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "User identifier.",
			},
			"passkey_status":   schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Passkey MFA status fields."},
			"email_mfa_status": schema.MapAttribute{ElementType: types.StringType, Computed: true, MarkdownDescription: "Email MFA status fields."},
			"passkey_credentials": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Registered passkey credentials.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":           schema.StringAttribute{Computed: true, MarkdownDescription: "Credential identifier."},
						"name":         schema.StringAttribute{Computed: true, MarkdownDescription: "Credential label."},
						"created_at":   schema.StringAttribute{Computed: true, MarkdownDescription: "Registration timestamp."},
						"last_used_at": schema.StringAttribute{Computed: true, MarkdownDescription: "Last usage timestamp."},
					},
				},
			},
			"email_mfa_codes": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Email MFA codes / purposes on file.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.StringAttribute{Computed: true, MarkdownDescription: "Code entry identifier."},
						"purpose":    schema.StringAttribute{Computed: true, MarkdownDescription: "Purpose of the code."},
						"status":     schema.StringAttribute{Computed: true, MarkdownDescription: "Code status."},
						"created_at": schema.StringAttribute{Computed: true, MarkdownDescription: "Creation timestamp."},
					},
				},
			},
		},
	}
}

func (d *mfaStatusDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *mfaStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data mfaStatusModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userType := data.UserType.ValueString()
	userID := data.UserID.ValueString()
	data.ID = types.StringValue(fmt.Sprintf("%s/%s", userType, userID))

	if status, err := d.client.GetPasskeyMFAStatus(ctx, userType, userID); err == nil {
		data.PasskeyStatus = common.MapStrings(client.StringMap(status))
	}
	if status, err := d.client.GetEmailMFAStatus(ctx, userType, userID); err == nil {
		data.EmailMFAStatus = common.MapStrings(client.StringMap(status))
	}
	if creds, err := d.client.ListPasskeyCredentials(ctx, userType, userID); err == nil {
		data.PasskeyCredentials = make([]passkeyRow, 0, len(creds))
		for _, item := range creds {
			m, ok := client.AsMap(item)
			if !ok {
				continue
			}
			data.PasskeyCredentials = append(data.PasskeyCredentials, passkeyRow{
				ID:         common.StringOrNull(client.FirstString(m, "id", "credential_id")),
				Name:       common.StringOrNull(client.FirstString(m, "name", "label", "title")),
				CreatedAt:  common.StringOrNull(client.FirstString(m, "created_at", "created")),
				LastUsedAt: common.StringOrNull(client.FirstString(m, "last_used_at", "last_used")),
			})
		}
	}
	if codes, err := d.client.ListEmailMFACodes(ctx, userType, userID); err == nil {
		data.EmailMFACodes = make([]mfaCodeRow, 0, len(codes))
		for _, item := range codes {
			m, ok := client.AsMap(item)
			if !ok {
				continue
			}
			data.EmailMFACodes = append(data.EmailMFACodes, mfaCodeRow{
				ID:        common.StringOrNull(client.FirstString(m, "id", "code_id")),
				Purpose:   common.StringOrNull(client.FirstString(m, "purpose", "type")),
				Status:    common.StringOrNull(client.FirstString(m, "status", "state")),
				CreatedAt: common.StringOrNull(client.FirstString(m, "created_at", "created")),
			})
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
