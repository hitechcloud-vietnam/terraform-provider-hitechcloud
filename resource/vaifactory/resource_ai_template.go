// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vaifactory

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
	_ resource.Resource                = (*aiTemplateResource)(nil)
	_ resource.ResourceWithConfigure   = (*aiTemplateResource)(nil)
	_ resource.ResourceWithImportState = (*aiTemplateResource)(nil)
)

// AITemplateResource returns the hitechcloud_ai_template resource.
func AITemplateResource() resource.Resource {
	return &aiTemplateResource{}
}

type aiTemplateResource struct {
	client *client.Client
}

type aiTemplateModel struct {
	ID          types.String `tfsdk:"id"`
	ServiceID   types.String `tfsdk:"service_id"`
	TemplateID  types.String `tfsdk:"template_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	IsPublic    types.Bool   `tfsdk:"is_public"`
}

func (r *aiTemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ai_template"
}

func (r *aiTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a HiTechCloud AI Factory instance template " +
			"(`POST /api/service/{service_id}/templates`). Name, description and " +
			"visibility can be updated in place.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `<service_id>/<template_id>`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the AI Factory service. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"template_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Template identifier.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Template name.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Template description.",
			},
			"is_public": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Whether the template is visible to other users.",
			},
		},
	}
}

func (r *aiTemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *aiTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aiTemplateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := client.AITemplate{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		IsPublic:    plan.IsPublic.ValueBool(),
	}
	templateID, err := r.client.CreateAITemplate(ctx, plan.ServiceID.ValueString(), in)
	if err != nil {
		resp.Diagnostics.AddError("Error creating template", err.Error())
		return
	}

	plan.TemplateID = types.StringValue(templateID)
	plan.ID = types.StringValue(client.FormatID(plan.ServiceID.ValueString(), templateID))
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *aiTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state aiTemplateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, templateID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}
	state.ServiceID = types.StringValue(serviceID)
	state.TemplateID = types.StringValue(templateID)
	if !r.refresh(ctx, &state, &resp.Diagnostics) {
		resp.State.RemoveResource(ctx)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *aiTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan aiTemplateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, templateID, ok := r.parseID(ctx, &plan, &resp.Diagnostics)
	if !ok {
		return
	}

	in := client.AITemplate{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		IsPublic:    plan.IsPublic.ValueBool(),
	}
	if err := r.client.UpdateAITemplate(ctx, serviceID, templateID, in); err != nil {
		resp.Diagnostics.AddError("Error updating template", err.Error())
		return
	}

	plan.ID = types.StringValue(client.FormatID(serviceID, templateID))
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *aiTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state aiTemplateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, templateID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	if err := r.client.DeleteAITemplate(ctx, serviceID, templateID); client.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error deleting template", err.Error())
	}
}

func (r *aiTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := client.SplitID(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("template_id"), parts[1])...)
}

func (r *aiTemplateResource) parseID(_ context.Context, state *aiTemplateModel, diags *diag.Diagnostics) (string, string, bool) {
	if state.ServiceID.ValueString() != "" && state.TemplateID.ValueString() != "" {
		return state.ServiceID.ValueString(), state.TemplateID.ValueString(), true
	}
	parts, err := client.SplitID(state.ID.ValueString(), 2)
	if err != nil {
		diags.AddError("Unexpected Resource ID", err.Error())
		return "", "", false
	}
	return parts[0], parts[1], true
}

// refresh reloads the template; returns false when it no longer exists.
func (r *aiTemplateResource) refresh(ctx context.Context, model *aiTemplateModel, diags *diag.Diagnostics) bool {
	tpl, err := r.client.GetAITemplate(ctx, model.ServiceID.ValueString(), model.TemplateID.ValueString())
	if client.IsNotFound(err) {
		return false
	}
	if err != nil {
		diags.AddError("Error reading template", err.Error())
		return true
	}
	model.Name = common.StringOrNull(tpl.Name)
	model.Description = common.StringOrNull(tpl.Description)
	model.IsPublic = types.BoolValue(tpl.IsPublic)
	return true
}
