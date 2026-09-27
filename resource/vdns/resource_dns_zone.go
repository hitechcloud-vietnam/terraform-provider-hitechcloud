// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vdns

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	_ resource.Resource                = (*dnsZoneResource)(nil)
	_ resource.ResourceWithConfigure   = (*dnsZoneResource)(nil)
	_ resource.ResourceWithImportState = (*dnsZoneResource)(nil)
)

// DNSZoneResource returns the hitechcloud_dns_zone resource.
func DNSZoneResource() resource.Resource {
	return &dnsZoneResource{}
}

type dnsZoneResource struct {
	client *client.Client
}

type dnsZoneModel struct {
	ID        types.String `tfsdk:"id"`
	ServiceID types.String `tfsdk:"service_id"`
	ZoneID    types.String `tfsdk:"zone_id"`
	Name      types.String `tfsdk:"name"`
	Records   types.List   `tfsdk:"records"`
}

// dnsRecordItem is the shared nested object for DNS records.
type dnsRecordItem struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Type     types.String `tfsdk:"type"`
	Content  types.String `tfsdk:"content"`
	TTL      types.Int64  `tfsdk:"ttl"`
	Priority types.Int64  `tfsdk:"priority"`
}

// dnsRecordSchema returns the nested attribute schema shared by data sources.
func dnsRecordSchema() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Record identifier inside the zone.",
			},
			"name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Record name (host).",
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Record type (A, AAAA, CNAME, MX, TXT, ...).",
			},
			"content": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Record value.",
			},
			"ttl": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Time to live in seconds.",
			},
			"priority": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Priority (MX/SRV records).",
			},
		},
	}
}

func (r *dnsZoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_zone"
}

func (r *dnsZoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a DNS zone on a HiTechCloud DNS service. " +
			"The zone is created under `POST /api/service/{service_id}/dns` and " +
			"removed with `DELETE /api/service/{service_id}/dns/{zone_id}`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `<service_id>/<zone_id>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the DNS service that owns the zone. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"zone_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier of the DNS zone.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Zone name (domain), e.g. `example.com`. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"records": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Records currently present in the zone (read-only).",
				NestedObject:        dnsRecordSchema(),
			},
		},
	}
}

func (r *dnsZoneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dnsZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dnsZoneModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneID, err := r.client.CreateDNSZone(ctx, plan.ServiceID.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating DNS zone", err.Error())
		return
	}

	plan.ZoneID = types.StringValue(zoneID)
	plan.ID = types.StringValue(client.FormatID(plan.ServiceID.ValueString(), zoneID))

	// Populate the computed records listing right after creation.
	if zone, err := r.client.GetDNSZone(ctx, plan.ServiceID.ValueString(), zoneID); err == nil {
		plan.Name = common.StringOrNull(zone.Name)
		plan.Records = recordsListValue(ctx, zone.Records)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dnsZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dnsZoneModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, zoneID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	zone, err := r.client.GetDNSZone(ctx, serviceID, zoneID)
	if client.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading DNS zone", err.Error())
		return
	}

	state.ServiceID = types.StringValue(serviceID)
	state.ZoneID = types.StringValue(zoneID)
	state.Name = common.StringOrNull(zone.Name)
	state.Records = recordsListValue(ctx, zone.Records)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *dnsZoneResource) Update(ctx context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All attributes force replacement; updates are not supported by the API.
	resp.Diagnostics.AddError(
		"Update not supported",
		"The HiTechCloud API provides no endpoint to rename or move a DNS zone; recreate the resource instead.",
	)
}

func (r *dnsZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dnsZoneModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, zoneID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	if err := r.client.DeleteDNSZone(ctx, serviceID, zoneID); client.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error deleting DNS zone", err.Error())
	}
}

func (r *dnsZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := client.SplitID(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_id"), parts[1])...)
}

// parseID resolves service and zone ids from the composite id / state.
func (r *dnsZoneResource) parseID(ctx context.Context, state *dnsZoneModel, diags *diag.Diagnostics) (string, string, bool) {
	_ = ctx
	if state.ServiceID.ValueString() != "" && state.ZoneID.ValueString() != "" {
		return state.ServiceID.ValueString(), state.ZoneID.ValueString(), true
	}
	parts, err := client.SplitID(state.ID.ValueString(), 2)
	if err != nil {
		diags.AddError("Unexpected Resource ID", err.Error())
		return "", "", false
	}
	return parts[0], parts[1], true
}

func toRecordItems(records []client.DNSRecord) []dnsRecordItem {
	items := make([]dnsRecordItem, 0, len(records))
	for _, rec := range records {
		items = append(items, dnsRecordItem{
			ID:       common.StringOrNull(rec.ID),
			Name:     common.StringOrNull(rec.Name),
			Type:     common.StringOrNull(rec.Type),
			Content:  common.StringOrNull(rec.Content),
			TTL:      common.Int64OrNull(rec.TTL),
			Priority: common.Int64OrNull(rec.Priority),
		})
	}
	return items
}

// dnsRecordItemAttrTypes mirrors the dnsRecordSchema nested attributes.
func dnsRecordItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":       types.StringType,
		"name":     types.StringType,
		"type":     types.StringType,
		"content":  types.StringType,
		"ttl":      types.Int64Type,
		"priority": types.Int64Type,
	}
}

// recordsListValue converts records into a types.List so computed values may
// stay unknown in plans (plain Go slices cannot hold unknown values).
func recordsListValue(ctx context.Context, records []client.DNSRecord) types.List {
	lv, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: dnsRecordItemAttrTypes()}, toRecordItems(records))
	if diags.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: dnsRecordItemAttrTypes()})
	}
	return lv
}
