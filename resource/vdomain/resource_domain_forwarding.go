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
	_ resource.Resource                = (*domainForwardingResource)(nil)
	_ resource.ResourceWithConfigure   = (*domainForwardingResource)(nil)
	_ resource.ResourceWithImportState = (*domainForwardingResource)(nil)
)

// DomainForwardingResource returns the hitechcloud_domain_forwarding resource.
func DomainForwardingResource() resource.Resource {
	return &domainForwardingResource{}
}

type domainForwardingResource struct {
	client *client.Client
}

type domainForwardingModel struct {
	ID       types.String `tfsdk:"id"`
	DomainID types.String `tfsdk:"domain_id"`
	URL      types.String `tfsdk:"url"`
	Frame    types.Bool   `tfsdk:"frame"`
}

func (r *domainForwardingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_forwarding"
}

func (r *domainForwardingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages URL forwarding (web redirect) of a registered domain " +
			"(`PUT /api/domain/{id}/forwarding`). Destroying the resource disables forwarding.",
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
			"url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Destination URL the domain should redirect to.",
			},
			"frame": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When true, the destination is rendered inside a frame instead of an HTTP redirect.",
			},
		},
	}
}

func (r *domainForwardingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *domainForwardingResource) apply(ctx context.Context, model *domainForwardingModel, clear bool) error {
	params := map[string]string{"url": model.URL.ValueString()}
	if clear {
		params["url"] = ""
	}
	if !model.Frame.IsNull() && !model.Frame.IsUnknown() {
		if model.Frame.ValueBool() {
			params["frame"] = "1"
		} else {
			params["frame"] = "0"
		}
	}
	return r.client.UpdateDomainForwarding(ctx, model.DomainID.ValueString(), params)
}

func (r *domainForwardingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan domainForwardingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.apply(ctx, &plan, false); err != nil {
		resp.Diagnostics.AddError("Error enabling domain forwarding", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.DomainID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *domainForwardingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state domainForwardingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The API exposes forwarding state through the domain detail; treat the
	// configured values as authoritative and only drop the resource when the
	// domain itself is gone.
	if _, err := r.client.GetDomainContactInfo(ctx, state.DomainID.ValueString()); err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *domainForwardingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan domainForwardingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.apply(ctx, &plan, false); err != nil {
		resp.Diagnostics.AddError("Error updating domain forwarding", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *domainForwardingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state domainForwardingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.apply(ctx, &state, true); err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error disabling domain forwarding", err.Error())
	}
}

func (r *domainForwardingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_id"), req.ID)...)
}
