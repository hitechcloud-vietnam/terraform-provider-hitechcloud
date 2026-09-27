// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vdns

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
	_ resource.Resource                = (*domainDNSRecordResource)(nil)
	_ resource.ResourceWithConfigure   = (*domainDNSRecordResource)(nil)
	_ resource.ResourceWithImportState = (*domainDNSRecordResource)(nil)
)

// DomainDNSRecordResource returns the hitechcloud_domain_dns_record resource.
func DomainDNSRecordResource() resource.Resource {
	return &domainDNSRecordResource{}
}

type domainDNSRecordResource struct {
	client *client.Client
}

type domainDNSRecordModel struct {
	ID       types.String `tfsdk:"id"`
	DomainID types.String `tfsdk:"domain_id"`
	RecordID types.String `tfsdk:"record_id"`
	Name     types.String `tfsdk:"name"`
	Type     types.String `tfsdk:"type"`
	Content  types.String `tfsdk:"content"`
	Priority types.Int64  `tfsdk:"priority"`
}

func (r *domainDNSRecordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_dns_record"
}

func (r *domainDNSRecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a DNS record of a registered domain (zone editor), " +
			"as configured through `POST /api/domain/{id}/dns`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `<domain_id>/<record_id>`.",
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
			"record_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Record index assigned by the API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Record name (host), e.g. `www` or `@`.",
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Record type (A, AAAA, CNAME, MX, TXT, ...).",
			},
			"content": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Record value.",
			},
			"priority": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Priority for MX/SRV records.",
			},
		},
	}
}

func (r *domainDNSRecordResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *domainDNSRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan domainDNSRecordModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rec := client.DomainDNSRecord{
		Name:     plan.Name.ValueString(),
		Type:     plan.Type.ValueString(),
		Content:  plan.Content.ValueString(),
		Priority: plan.Priority.ValueInt64(),
	}
	recordID, err := r.client.CreateDomainDNSRecord(ctx, plan.DomainID.ValueString(), rec)
	if err != nil {
		resp.Diagnostics.AddError("Error creating domain DNS record", err.Error())
		return
	}

	plan.RecordID = types.StringValue(recordID)
	plan.ID = types.StringValue(client.FormatID(plan.DomainID.ValueString(), recordID))
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *domainDNSRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state domainDNSRecordModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainID, recordID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	state.DomainID = types.StringValue(domainID)
	state.RecordID = types.StringValue(recordID)
	if !r.refresh(ctx, &state, &resp.Diagnostics) {
		resp.State.RemoveResource(ctx)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// refresh reloads the record; returns false when it no longer exists.
func (r *domainDNSRecordResource) refresh(ctx context.Context, model *domainDNSRecordModel, diags *diag.Diagnostics) bool {
	rec, err := r.client.GetDomainDNSRecord(ctx, model.DomainID.ValueString(), model.RecordID.ValueString())
	if client.IsNotFound(err) {
		return false
	}
	if err != nil {
		diags.AddError("Error reading domain DNS record", err.Error())
		return true
	}
	model.Name = common.StringOrNull(rec.Name)
	model.Type = common.StringOrNull(rec.Type)
	model.Content = common.StringOrNull(rec.Content)
	model.Priority = common.Int64OrNull(rec.Priority)
	return true
}

func (r *domainDNSRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan domainDNSRecordModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainID, recordID, ok := r.parseID(ctx, &plan, &resp.Diagnostics)
	if !ok {
		return
	}

	rec := client.DomainDNSRecord{
		Name:     plan.Name.ValueString(),
		Type:     plan.Type.ValueString(),
		Content:  plan.Content.ValueString(),
		Priority: plan.Priority.ValueInt64(),
	}
	if err := r.client.UpdateDomainDNSRecord(ctx, domainID, recordID, rec); err != nil {
		resp.Diagnostics.AddError("Error updating domain DNS record", err.Error())
		return
	}

	plan.ID = types.StringValue(client.FormatID(domainID, recordID))
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *domainDNSRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state domainDNSRecordModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainID, recordID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	if err := r.client.DeleteDomainDNSRecord(ctx, domainID, recordID); client.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error deleting domain DNS record", err.Error())
	}
}

func (r *domainDNSRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := client.SplitID(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("record_id"), parts[1])...)
}

func (r *domainDNSRecordResource) parseID(_ context.Context, state *domainDNSRecordModel, diags *diag.Diagnostics) (string, string, bool) {
	if state.DomainID.ValueString() != "" && state.RecordID.ValueString() != "" {
		return state.DomainID.ValueString(), state.RecordID.ValueString(), true
	}
	parts, err := client.SplitID(state.ID.ValueString(), 2)
	if err != nil {
		diags.AddError("Unexpected Resource ID", err.Error())
		return "", "", false
	}
	return parts[0], parts[1], true
}
