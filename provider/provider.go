// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

// Package provider implements the Terraform provider for the HiTechCloud User
// API. The repository layout follows github.com/vngcloud/terraform-provider-vngcloud:
// the HTTP client lives in /client, the provider entry point in /provider and
// the Terraform resources and data sources in /resource/<group>.
package provider

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/client"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/vaccount"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/vaifactory"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/vbilling"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/vdns"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/vdomain"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/vportal"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/vserver"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/vstorage"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/vsupport"
)

// Environment variable names honoured by the provider.
const (
	EnvToken    = "HITECHCLOUD_TOKEN"
	EnvEndpoint = "HITECHCLOUD_ENDPOINT"
	EnvTimeout  = "HITECHCLOUD_REQUEST_TIMEOUT"
	EnvUsername = "HITECHCLOUD_USERNAME"
	EnvPassword = "HITECHCLOUD_PASSWORD"
	EnvRefresh  = "HITECHCLOUD_REFRESH_TOKEN"
)

// Ensure the provider fully satisfies the framework interfaces.
var _ provider.Provider = (*hiTechCloudProvider)(nil)

// hiTechCloudProvider implements provider.Provider.
type hiTechCloudProvider struct {
	version string
}

// New returns the provider factory used by main.go.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &hiTechCloudProvider{version: version}
	}
}

// providerModel maps provider schema data.
type providerModel struct {
	Token        types.String `tfsdk:"token"`
	RefreshToken types.String `tfsdk:"refresh_token"`
	Username     types.String `tfsdk:"username"`
	Password     types.String `tfsdk:"password"`
	Endpoint     types.String `tfsdk:"endpoint"`
	Timeout      types.String `tfsdk:"request_timeout"`
}

func (p *hiTechCloudProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hitechcloud"
	resp.Version = p.version
}

func (p *hiTechCloudProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The HiTechCloud provider manages resources exposed by the " +
			"[HiTechCloud User API](https://my.hitechcloud.vn) (`https://api.hitechcloud.vn`). " +
			"Authentication uses the two tokens returned by `POST /api/login`: an access token " +
			"(`token`) and a refresh token (`refresh_token`). Provide either `token` directly, " +
			"or `username`/`password` to log in automatically; the refresh token renews the " +
			"access token whenever it expires.",
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				MarkdownDescription: "HiTechCloud API access token (the `token` field of the `POST /api/login` response). " +
					"May also be provided via the `" + EnvToken + "` environment variable. " +
					"The token is marked sensitive and is never logged. " +
					"Not needed when `username`/`password` are set.",
				Optional:  true,
				Sensitive: true,
			},
			"refresh_token": schema.StringAttribute{
				MarkdownDescription: "HiTechCloud API refresh token (the `refresh_token` field of the `POST /api/login` response). " +
					"Used by `POST /api/token` to obtain a fresh access token whenever the current one expires. " +
					"May also be provided via the `" + EnvRefresh + "` environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Account email address used for `POST /api/login`. " +
					"May also be provided via the `" + EnvUsername + "` environment variable.",
				Optional: true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Account password used for `POST /api/login`. " +
					"May also be provided via the `" + EnvPassword + "` environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Base URL of the HiTechCloud User API. " +
					"Defaults to `https://api.hitechcloud.vn` (the `baseUrl` variable of the " +
					"API specification). May also be provided via the `" + EnvEndpoint + "` environment variable.",
				Optional: true,
			},
			"request_timeout": schema.StringAttribute{
				MarkdownDescription: "Per-request timeout as a Go duration string (e.g. `60s`, `2m`). " +
					"Defaults to `60s`. May also be provided via the `" + EnvTimeout + "` environment variable.",
				Optional: true,
			},
		},
	}
}

func (p *hiTechCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	token := os.Getenv(EnvToken)
	if !data.Token.IsNull() && !data.Token.IsUnknown() && data.Token.ValueString() != "" {
		token = data.Token.ValueString()
	}

	refreshToken := os.Getenv(EnvRefresh)
	if !data.RefreshToken.IsNull() && !data.RefreshToken.IsUnknown() && data.RefreshToken.ValueString() != "" {
		refreshToken = data.RefreshToken.ValueString()
	}

	username := os.Getenv(EnvUsername)
	if !data.Username.IsNull() && !data.Username.IsUnknown() && data.Username.ValueString() != "" {
		username = data.Username.ValueString()
	}

	password := os.Getenv(EnvPassword)
	if !data.Password.IsNull() && !data.Password.IsUnknown() && data.Password.ValueString() != "" {
		password = data.Password.ValueString()
	}

	endpoint := os.Getenv(EnvEndpoint)
	if !data.Endpoint.IsNull() && !data.Endpoint.IsUnknown() && data.Endpoint.ValueString() != "" {
		endpoint = data.Endpoint.ValueString()
	}
	if endpoint == "" {
		endpoint = client.DefaultEndpoint
	}

	timeout := 60 * time.Second
	if v := os.Getenv(EnvTimeout); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("request_timeout"),
				"Invalid "+EnvTimeout,
				fmt.Sprintf("Could not parse %q as duration: %s", v, err),
			)
			return
		}
		timeout = d
	}
	if !data.Timeout.IsNull() && !data.Timeout.IsUnknown() && data.Timeout.ValueString() != "" {
		d, err := time.ParseDuration(data.Timeout.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("request_timeout"),
				"Invalid request_timeout",
				fmt.Sprintf("Could not parse %q as duration: %s", data.Timeout.ValueString(), err),
			)
			return
		}
		timeout = d
	}

	cli, err := client.New(client.Config{
		Endpoint:     endpoint,
		Token:        token,
		RefreshToken: refreshToken,
		Timeout:      timeout,
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to create HiTechCloud client", err.Error())
		return
	}

	// No static token: log in with username/password. POST /api/login returns
	// both tokens ("đăng nhập sinh 2 token"): token and refresh_token.
	if token == "" {
		if username == "" || password == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("token"),
				"Missing HiTechCloud API credentials",
				"Provide a token ("+EnvToken+"), or username and password ("+EnvUsername+
					"/"+EnvPassword+") to log in via POST /api/login. "+
					"An optional refresh_token ("+EnvRefresh+") renews the access token automatically.",
			)
			return
		}
		login, err := cli.Login(ctx, username, password)
		if err != nil {
			resp.Diagnostics.AddError("HiTechCloud login failed", err.Error())
			return
		}
		if login.RefreshToken != "" {
			cli.SetRefreshToken(login.RefreshToken)
		}
	}

	resp.ResourceData = cli
	resp.DataSourceData = cli
}

func (p *hiTechCloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		vaccount.AccountDataSource,
		vaccount.ContactsDataSource,
		vaccount.AccountLogsDataSource,
		vserver.ServicesDataSource,
		vserver.ServiceDataSource,
		vserver.PVEDataSource,
		vserver.IPAMDataSource,
		vbilling.InvoicesDataSource,
		vbilling.PaymentMethodsDataSource,
		vbilling.ProductsDataSource,
		vbilling.CertificatesDataSource,
		vbilling.WillExpiredDataSource,
		vdomain.DomainsDataSource,
		vdomain.DomainDataSource,
		vdomain.DomainTLDsDataSource,
		vdomain.DNSSECKeysDataSource,
		vdomain.DomainContactDataSource,
		vdns.DNSZonesDataSource,
		vserver.VMsDataSource,
		vaifactory.AIInstanceTypesDataSource,
		vaifactory.AISSHKeysDataSource,
		vaifactory.AIVolumesDataSource,
		vaifactory.AIClusterTypesDataSource,
		vstorage.S3BucketsDataSource,
		vstorage.S3SubusersDataSource,
		vstorage.PBSDataSource,
		vstorage.S3ConnectionDataSource,
		vbilling.BalanceDataSource,
		vbilling.CategoriesDataSource,
		vbilling.PaymentFeesDataSource,
		vdomain.WhoisDataSource,
		vdomain.DomainDNSTypesDataSource,
		vdomain.DomainAvailabilityDataSource,
		vportal.URLShortenerLinksDataSource,
		vportal.PartnerDataSource,
		vsupport.TicketsDataSource,
		vsupport.TicketDepartmentsDataSource,
		vsupport.NotificationsDataSource,
		vsupport.StatusesDataSource,
	}
}

func (p *hiTechCloudProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		vdns.DNSZoneResource,
		vdns.DNSRecordResource,
		vdns.DomainDNSRecordResource,
		vdomain.DomainSettingsResource,
		vserver.VMResource,
		vserver.VMInterfaceResource,
		vserver.VMFirewallRuleResource,
		vserver.ServiceIPResource,
		vaifactory.AIInstanceResource,
		vaifactory.AISSHKeyResource,
		vaifactory.AIVolumeResource,
		vaifactory.AITemplateResource,
		vaifactory.AIClusterResource,
		vstorage.S3BucketResource,
		vstorage.S3SubuserResource,
		vportal.URLShortenerLinkResource,
		vaccount.ContactResource,
		vsupport.TicketResource,
		vserver.RDNSResource,
		vdomain.DNSSECKeyResource,
		vdomain.EmailForwardingResource,
		vdomain.DomainForwardingResource,
		vbilling.ServiceAutoRenewResource,
		vserver.VMSnapshotResource,
		vportal.PartnerLeadResource,
	}
}
