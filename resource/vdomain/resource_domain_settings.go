// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vdomain

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/client"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/common"
)

var (
	_ resource.Resource                = (*domainSettingsResource)(nil)
	_ resource.ResourceWithConfigure   = (*domainSettingsResource)(nil)
	_ resource.ResourceWithImportState = (*domainSettingsResource)(nil)
)

// DomainSettingsResource returns the hitechcloud_domain_settings resource.
func DomainSettingsResource() resource.Resource {
	return &domainSettingsResource{}
}

type domainSettingsResource struct {
	client *client.Client
}

type domainSettingsModel struct {
	ID            types.String `tfsdk:"id"`
	DomainID      types.String `tfsdk:"domain_id"`
	Nameservers   types.Set    `tfsdk:"nameservers"`
	Autorenew     types.Bool   `tfsdk:"autorenew"`
	RegistrarLock types.Bool   `tfsdk:"registrar_lock"`
	IDProtection  types.Bool   `tfsdk:"id_protection"`
}

func (r *domainSettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_settings"
}

func (r *domainSettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the settings of a registered domain: nameservers, " +
			"auto-renew, registrar lock and ID protection " +
			"(`PUT /api/domain/{id}/ns|autorenew|reglock|idprotection`).\n\n" +
			"Destroying this resource only removes it from Terraform state; domain " +
			"settings themselves are kept as-is because the API has no reset endpoint.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The domain identifier (same as `domain_id`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Identifier of the registered domain. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"nameservers": schema.SetAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Nameservers of the domain. An empty set restores the default nameservers.",
			},
			"autorenew": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the domain renews automatically.",
			},
			"registrar_lock": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether the registrar (transfer) lock is enabled.",
			},
			"id_protection": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether WHOIS ID protection is enabled.",
			},
		},
	}
}

func (r *domainSettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	cli, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got %T", req.ProviderData),
		)
		return
	}
	r.client = cli
}

func (r *domainSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan domainSettingsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(plan.DomainID.ValueString())
	if diags := r.apply(ctx, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *domainSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state domainSettingsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.DomainID.IsNull() || state.DomainID.ValueString() == "" {
		parts, err := client.SplitID(state.ID.ValueString(), 1)
		if err != nil {
			resp.Diagnostics.AddError("Unexpected Resource ID", err.Error())
			return
		}
		state.DomainID = types.StringValue(parts[0])
	}
	if !r.refresh(ctx, &state, &resp.Diagnostics) {
		resp.State.RemoveResource(ctx)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *domainSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan domainSettingsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if diags := r.apply(ctx, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete is intentionally a no-op: the API provides no endpoint to reset
// domain settings, so the resource is only removed from Terraform state.
func (r *domainSettingsResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

func (r *domainSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_id"), req.ID)...)
}

// apply pushes the configured settings to the API.
func (r *domainSettingsResource) apply(ctx context.Context, plan *domainSettingsModel) diag.Diagnostics {
	var diags diag.Diagnostics
	domainID := plan.DomainID.ValueString()

	if !plan.Nameservers.IsNull() && !plan.Nameservers.IsUnknown() {
		ns, ds := common.StringsFromSet(ctx, plan.Nameservers)
		diags.Append(ds...)
		if err := r.client.SetDomainNameservers(ctx, domainID, ns); err != nil {
			diags.AddError("Error setting nameservers", err.Error())
		}
	}
	if !plan.Autorenew.IsNull() && !plan.Autorenew.IsUnknown() {
		if err := r.client.SetDomainAutorenew(ctx, domainID, plan.Autorenew.ValueBool()); err != nil {
			diags.AddError("Error setting auto-renew", err.Error())
		}
	}
	if !plan.RegistrarLock.IsNull() && !plan.RegistrarLock.IsUnknown() {
		if err := r.client.SetDomainReglock(ctx, domainID, plan.RegistrarLock.ValueBool()); err != nil {
			diags.AddError("Error setting registrar lock", err.Error())
		}
	}
	if !plan.IDProtection.IsNull() && !plan.IDProtection.IsUnknown() {
		if err := r.client.SetDomainIDProtection(ctx, domainID, plan.IDProtection.ValueBool()); err != nil {
			diags.AddError("Error setting ID protection", err.Error())
		}
	}
	return diags
}

// refresh reads the current settings from the API into the model. It returns
// false when the domain no longer exists so Read can drop the resource.
func (r *domainSettingsResource) refresh(ctx context.Context, model *domainSettingsModel, diags *diag.Diagnostics) bool {
	domain, err := r.client.GetDomain(ctx, model.DomainID.ValueString())
	if client.IsNotFound(err) {
		return false
	}
	if err != nil {
		diags.AddError("Error reading domain settings", err.Error())
		return true
	}
	model.ID = types.StringValue(model.DomainID.ValueString())
	model.Nameservers = common.SetStrings(domain.Nameservers)
	model.Autorenew = types.BoolValue(domain.Autorenew)
	model.RegistrarLock = types.BoolValue(domain.RegistrarLock)
	model.IDProtection = types.BoolValue(domain.IDProtection)
	return true
}

// splitID keeps parity with other packages (single-part id).
func splitID(id string) (string, error) {
	parts, err := client.SplitID(id, 1)
	if err != nil {
		return "", common.FormatIDError(id, 1)
	}
	return parts[0], nil
}
