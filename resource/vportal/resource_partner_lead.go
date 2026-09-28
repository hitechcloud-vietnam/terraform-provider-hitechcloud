// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vportal

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
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/common"
)

var (
	_ resource.Resource                = (*partnerLeadResource)(nil)
	_ resource.ResourceWithConfigure   = (*partnerLeadResource)(nil)
	_ resource.ResourceWithImportState = (*partnerLeadResource)(nil)
)

// PartnerLeadResource returns the hitechcloud_partner_lead resource.
func PartnerLeadResource() resource.Resource {
	return &partnerLeadResource{}
}

type partnerLeadResource struct {
	client *client.Client
}

type partnerLeadModel struct {
	ID          types.String `tfsdk:"id"`
	Email       types.String `tfsdk:"email"`
	Company     types.String `tfsdk:"company"`
	ContactName types.String `tfsdk:"contact_name"`
	Phone       types.String `tfsdk:"phone"`
	Status      types.String `tfsdk:"status"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

func (r *partnerLeadResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_partner_lead"
}

func (r *partnerLeadResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Registers a partner lead / sales opportunity " +
			"(`POST /api/partner/leads`). Leads are immutable registrations: " +
			"changing any attribute forces a new resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Lead identifier returned by the API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Lead's email address; duplicates with an active lead are rejected by the API. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"company": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Company name of the lead. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"contact_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Contact person's name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"phone": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Contact phone number.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Current lead status reported by the API.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Creation timestamp reported by the API.",
			},
		},
	}
}

func (r *partnerLeadResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *partnerLeadResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan partnerLeadModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := map[string]string{
		"email":   plan.Email.ValueString(),
		"company": plan.Company.ValueString(),
	}
	if v := plan.ContactName.ValueString(); v != "" {
		params["contact_name"] = v
	}
	if v := plan.Phone.ValueString(); v != "" {
		params["phone"] = v
	}

	lead, err := r.client.RegisterPartnerLead(ctx, params)
	if err != nil {
		resp.Diagnostics.AddError("Error registering partner lead", err.Error())
		return
	}

	id := client.FirstString(lead, "id", "lead_id", "hash")
	if id == "" {
		id = plan.Email.ValueString()
	}
	plan.ID = types.StringValue(id)
	plan.Status = common.StringOrNull(client.FirstString(lead, "status"))
	plan.CreatedAt = common.StringOrNull(client.FirstString(lead, "created_at", "created"))
	if v := client.FirstString(lead, "contact_name"); v != "" {
		plan.ContactName = types.StringValue(v)
	}
	if v := client.FirstString(lead, "phone"); v != "" {
		plan.Phone = types.StringValue(v)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *partnerLeadResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state partnerLeadModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	leads, err := r.client.ListPartnerLeads(ctx, map[string]string{})
	if err != nil {
		resp.Diagnostics.AddError("Error reading partner leads", err.Error())
		return
	}

	m, ok := client.FindByID(leads, state.ID.ValueString(), "id", "lead_id")
	if !ok {
		// Fall back to matching by email for older registrations.
		for _, item := range leads {
			mm, mapOK := client.AsMap(item)
			if mapOK && client.FirstString(mm, "email") == state.Email.ValueString() {
				m, ok = mm, true
				break
			}
		}
	}
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}
	state.Status = common.StringOrNull(client.FirstString(m, "status"))
	state.CreatedAt = common.StringOrNull(client.FirstString(m, "created_at", "created"))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Leads are immutable registrations: Update is only reached via replacement.
func (r *partnerLeadResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Partner leads are immutable",
		"The API does not support editing leads; this resource only supports create and read.",
	)
}

func (r *partnerLeadResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"Lead not deleted",
		"The HiTechCloud API does not expose a lead deletion endpoint; the lead was removed from Terraform state only.",
	)
}

func (r *partnerLeadResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("email"), req.ID)...)
}
