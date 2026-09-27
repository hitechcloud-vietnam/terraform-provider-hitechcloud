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
	Token    types.String `tfsdk:"token"`
	Endpoint types.String `tfsdk:"endpoint"`
	Timeout  types.String `tfsdk:"request_timeout"`
}

func (p *hiTechCloudProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hitechcloud"
	resp.Version = p.version
}

func (p *hiTechCloudProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The HiTechCloud provider manages resources exposed by the " +
			"[HiTechCloud User API](https://my.hitechcloud.vn) (`https://api.hitechcloud.vn`). " +
			"Authentication uses a bearer token obtained from `POST /api/login`; " +
			"see the provider documentation for how to obtain one.",
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				MarkdownDescription: "HiTechCloud API bearer token. " +
					"May also be provided via the `" + EnvToken + "` environment variable. " +
					"The token is marked sensitive and is never logged.",
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

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing HiTechCloud API token",
			"Set the token attribute or the "+EnvToken+" environment variable. "+
				"A token can be obtained by calling POST /api/login with your account credentials.",
		)
		return
	}

	cli, err := client.New(client.Config{
		Endpoint: endpoint,
		Token:    token,
		Timeout:  timeout,
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to create HiTechCloud client", err.Error())
		return
	}

	resp.ResourceData = cli
	resp.DataSourceData = cli
}

func (p *hiTechCloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		vaccount.AccountDataSource,
		vaccount.ContactsDataSource,
		vserver.ServicesDataSource,
		vserver.ServiceDataSource,
		vbilling.InvoicesDataSource,
		vbilling.PaymentMethodsDataSource,
		vbilling.ProductsDataSource,
		vbilling.CertificatesDataSource,
		vdomain.DomainsDataSource,
		vdomain.DomainDataSource,
		vdomain.DomainTLDsDataSource,
		vdns.DNSZonesDataSource,
		vserver.VMsDataSource,
		vaifactory.AIInstanceTypesDataSource,
		vaifactory.AISSHKeysDataSource,
		vaifactory.AIVolumesDataSource,
		vaifactory.AIClusterTypesDataSource,
		vstorage.S3BucketsDataSource,
		vstorage.S3SubusersDataSource,
		vbilling.BalanceDataSource,
		vbilling.CategoriesDataSource,
		vbilling.PaymentFeesDataSource,
		vdomain.WhoisDataSource,
		vdomain.DomainDNSTypesDataSource,
		vdomain.DomainAvailabilityDataSource,
		vportal.URLShortenerLinksDataSource,
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
	}
}
