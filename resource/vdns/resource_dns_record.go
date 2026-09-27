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
	_ resource.Resource                = (*dnsRecordResource)(nil)
	_ resource.ResourceWithConfigure   = (*dnsRecordResource)(nil)
	_ resource.ResourceWithImportState = (*dnsRecordResource)(nil)
)

// DNSRecordResource returns the hitechcloud_dns_record resource.
func DNSRecordResource() resource.Resource {
	return &dnsRecordResource{}
}

type dnsRecordResource struct {
	client *client.Client
}

type dnsRecordModel struct {
	ID        types.String `tfsdk:"id"`
	ServiceID types.String `tfsdk:"service_id"`
	ZoneID    types.String `tfsdk:"zone_id"`
	RecordID  types.String `tfsdk:"record_id"`
	Name      types.String `tfsdk:"name"`
	Type      types.String `tfsdk:"type"`
	Content   types.String `tfsdk:"content"`
	TTL       types.Int64  `tfsdk:"ttl"`
	Priority  types.Int64  `tfsdk:"priority"`
}

func (r *dnsRecordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

func (r *dnsRecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a record inside a HiTechCloud service DNS zone " +
			"(`POST /api/service/{service_id}/dns/{zone_id}/records`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `<service_id>/<zone_id>/<record_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the DNS service. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"zone_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "DNS zone identifier. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"record_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Record identifier assigned by the API.",
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
				MarkdownDescription: "Record value, e.g. an IP address or target host.",
			},
			"ttl": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Time to live in seconds.",
			},
			"priority": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Priority for MX/SRV records.",
			},
		},
	}
}

func (r *dnsRecordResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dnsRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dnsRecordModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rec := client.DNSRecord{
		Name:     plan.Name.ValueString(),
		Type:     plan.Type.ValueString(),
		Content:  plan.Content.ValueString(),
		TTL:      plan.TTL.ValueInt64(),
		Priority: plan.Priority.ValueInt64(),
	}
	recordID, err := r.client.CreateDNSRecord(ctx, plan.ServiceID.ValueString(), plan.ZoneID.ValueString(), rec)
	if err != nil {
		resp.Diagnostics.AddError("Error creating DNS record", err.Error())
		return
	}

	plan.RecordID = types.StringValue(recordID)
	plan.ID = types.StringValue(client.FormatID(plan.ServiceID.ValueString(), plan.ZoneID.ValueString(), recordID))
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dnsRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dnsRecordModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, zoneID, recordID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	state.ServiceID = types.StringValue(serviceID)
	state.ZoneID = types.StringValue(zoneID)
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

func (r *dnsRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dnsRecordModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, zoneID, recordID, ok := r.parseID(ctx, &plan, &resp.Diagnostics)
	if !ok {
		return
	}

	rec := client.DNSRecord{
		Name:     plan.Name.ValueString(),
		Type:     plan.Type.ValueString(),
		Content:  plan.Content.ValueString(),
		TTL:      plan.TTL.ValueInt64(),
		Priority: plan.Priority.ValueInt64(),
	}
	if err := r.client.UpdateDNSRecord(ctx, serviceID, zoneID, recordID, rec); err != nil {
		resp.Diagnostics.AddError("Error updating DNS record", err.Error())
		return
	}

	plan.ID = types.StringValue(client.FormatID(serviceID, zoneID, recordID))
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dnsRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dnsRecordModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, zoneID, recordID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	if err := r.client.DeleteDNSRecord(ctx, serviceID, zoneID, recordID); client.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error deleting DNS record", err.Error())
	}
}

func (r *dnsRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := client.SplitID(req.ID, 3)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_id"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("record_id"), parts[2])...)
}

func (r *dnsRecordResource) parseID(_ context.Context, state *dnsRecordModel, diags *diag.Diagnostics) (string, string, string, bool) {
	if state.ServiceID.ValueString() != "" && state.ZoneID.ValueString() != "" && state.RecordID.ValueString() != "" {
		return state.ServiceID.ValueString(), state.ZoneID.ValueString(), state.RecordID.ValueString(), true
	}
	parts, err := client.SplitID(state.ID.ValueString(), 3)
	if err != nil {
		diags.AddError("Unexpected Resource ID", err.Error())
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

// refresh reloads the record from the API. It returns false when the record no
// longer exists so Read can remove the resource from state.
func (r *dnsRecordResource) refresh(ctx context.Context, model *dnsRecordModel, diags *diag.Diagnostics) bool {
	rec, err := r.client.FindDNSRecord(ctx, model.ServiceID.ValueString(), model.ZoneID.ValueString(), model.RecordID.ValueString())
	if client.IsNotFound(err) {
		return false
	}
	if err != nil {
		diags.AddError("Error reading DNS record", err.Error())
		return true
	}
	model.Name = common.StringOrNull(rec.Name)
	model.Type = common.StringOrNull(rec.Type)
	model.Content = common.StringOrNull(rec.Content)
	model.TTL = common.Int64OrNull(rec.TTL)
	model.Priority = common.Int64OrNull(rec.Priority)
	return true
}
