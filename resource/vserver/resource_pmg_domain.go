// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vserver

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
	_ resource.Resource                = (*pmgDomainResource)(nil)
	_ resource.ResourceWithConfigure   = (*pmgDomainResource)(nil)
	_ resource.ResourceWithImportState = (*pmgDomainResource)(nil)
)

// PMGDomainResource returns the hitechcloud_pmg_domain resource.
func PMGDomainResource() resource.Resource {
	return &pmgDomainResource{}
}

type pmgDomainResource struct {
	client *client.Client
}

type pmgDomainModel struct {
	ID        types.String `tfsdk:"id"`
	ServiceID types.String `tfsdk:"service_id"`
	Domain    types.String `tfsdk:"domain"`
	Host      types.String `tfsdk:"host"`
	Port      types.String `tfsdk:"port"`
}

func (r *pmgDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pmg_domain"
}

func (r *pmgDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Protects a mail domain through a HiTechCloud PMG (Proxmox Mail " +
			"Gateway) service and sets its target mail transport " +
			"(`POST /api/service/{id}/htcpmg/domains`, `POST /api/service/{id}/htcpmg/transport`). " +
			"The API has no domain removal endpoint: destroying the resource removes it " +
			"from Terraform state only.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `service_id/domain`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HiTechCloud service ID of the PMG service. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"domain": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Mail domain to protect. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"host": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Target mail server (hostname or IP) receiving filtered mail.",
			},
			"port": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Target SMTP port, defaults to `25`.",
			},
		},
	}
}

func (r *pmgDomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *pmgDomainResource) applyTransport(ctx context.Context, model *pmgDomainModel) error {
	if model.Host.IsNull() || model.Host.ValueString() == "" {
		return nil
	}
	port := model.Port.ValueString()
	if port == "" {
		port = "25"
	}
	return r.client.SetPMGTransport(ctx,
		model.ServiceID.ValueString(),
		model.Domain.ValueString(),
		model.Host.ValueString(),
		port,
	)
}

func (r *pmgDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan pmgDomainModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.AddPMGDomain(ctx, plan.ServiceID.ValueString(), plan.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error adding PMG domain", err.Error())
		return
	}
	if err := r.applyTransport(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error setting PMG transport", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.ServiceID.ValueString(), plan.Domain.ValueString()))
	r.refresh(ctx, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *pmgDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state pmgDomainModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !r.refresh(ctx, &state) {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *pmgDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan pmgDomainModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.applyTransport(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Error updating PMG transport", err.Error())
		return
	}

	r.refresh(ctx, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *pmgDomainResource) Delete(_ context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		"PMG domain not removed",
		"The HiTechCloud API has no PMG domain removal endpoint; the domain was removed from Terraform state only.",
	)
}

func (r *pmgDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	parts, err := client.SplitID(id, 2)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf(
			"Expected import ID in the format service_id/domain: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), parts[1])...)
}

// refresh reloads the PMG config and reports whether the domain is protected.
func (r *pmgDomainResource) refresh(ctx context.Context, model *pmgDomainModel) bool {
	cfg, err := r.client.GetPMGConfig(ctx, model.ServiceID.ValueString())
	if err != nil {
		return !client.IsNotFound(err)
	}
	domains := client.ExtractList(cfg, "domains")
	for _, item := range domains {
		m, ok := client.AsMap(item)
		if !ok {
			continue
		}
		name := client.FirstString(m, "domain", "name")
		if name == model.Domain.ValueString() {
			model.Host = common.StringOrNull(client.FirstString(m, "host", "target", "relay"))
			model.Port = common.StringOrNull(client.FirstString(m, "port"))
			return true
		}
	}
	return false
}
