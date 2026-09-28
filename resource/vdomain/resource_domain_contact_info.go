// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vdomain

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/client"
)

var (
	_ resource.Resource                = (*domainContactResource)(nil)
	_ resource.ResourceWithConfigure   = (*domainContactResource)(nil)
	_ resource.ResourceWithImportState = (*domainContactResource)(nil)
)

// DomainContactResource returns the hitechcloud_domain_contact_info resource.
func DomainContactResource() resource.Resource {
	return &domainContactResource{}
}

type domainContactResource struct {
	client *client.Client
}

type domainContactResModel struct {
	ID          types.String `tfsdk:"id"`
	DomainID    types.String `tfsdk:"domain_id"`
	ContactInfo types.String `tfsdk:"contact_info"`
}

func (r *domainContactResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_contact_info"
}

func (r *domainContactResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the contact information of a registered domain " +
			"(`PUT /api/domain/{id}/contact`). The `contact_info` value is passed " +
			"verbatim to the API as the `contact_info` query parameter. Destroying " +
			"the resource leaves the contacts unchanged (the API has no reset).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The domain ID (same as `domain_id`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the registered domain. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"contact_info": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Contact information payload accepted by `PUT /api/domain/{id}/contact` (the `contact_info` parameter).",
			},
		},
	}
}

func (r *domainContactResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *domainContactResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan domainContactResModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateDomainContactInfo(ctx,
		plan.DomainID.ValueString(),
		plan.ContactInfo.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("Error setting domain contact info", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.DomainID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *domainContactResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state domainContactResModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.GetDomainContactInfo(ctx, state.DomainID.ValueString()); err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *domainContactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan domainContactResModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateDomainContactInfo(ctx,
		plan.DomainID.ValueString(),
		plan.ContactInfo.ValueString(),
	); err != nil {
		resp.Diagnostics.AddError("Error updating domain contact info", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *domainContactResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Domain contacts unchanged",
		"The HiTechCloud API has no contact reset; the attribute was removed from Terraform state only.",
	)
}

func (r *domainContactResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("contact_info"), "")...)
}
