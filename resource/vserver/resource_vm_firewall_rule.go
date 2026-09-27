// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vserver

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
	_ resource.Resource                = (*vmFirewallRuleResource)(nil)
	_ resource.ResourceWithConfigure   = (*vmFirewallRuleResource)(nil)
	_ resource.ResourceWithImportState = (*vmFirewallRuleResource)(nil)
)

// VMFirewallRuleResource returns the hitechcloud_vm_firewall_rule resource.
func VMFirewallRuleResource() resource.Resource {
	return &vmFirewallRuleResource{}
}

type vmFirewallRuleResource struct {
	client *client.Client
}

type vmFirewallRuleModel struct {
	ID           types.String `tfsdk:"id"`
	ServiceID    types.String `tfsdk:"service_id"`
	Position     types.Int64  `tfsdk:"position"`
	Action       types.String `tfsdk:"action"`
	Type         types.String `tfsdk:"type"`
	Comment      types.String `tfsdk:"comment"`
	Protocol     types.String `tfsdk:"protocol"`
	AddressStart types.String `tfsdk:"address_start"`
	AddressEnd   types.String `tfsdk:"address_end"`
	PortStart    types.Int64  `tfsdk:"port_start"`
	PortEnd      types.Int64  `tfsdk:"port_end"`
}

func (r *vmFirewallRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm_firewall_rule"
}

func (r *vmFirewallRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a virtual server firewall rule " +
			"(`POST /api/service/{service_id}/vms/firewall`). The API exposes no " +
			"rule update endpoint, so every change forces replacement.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `<service_id>/<position>`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the cloud service. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"position": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Rule position in the firewall list (assigned by the API).",
			},
			"action": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Rule action: `accept` or `drop`. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Rule type: `source`, `destination` or `both`. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Comment for the rule. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protocol": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Protocol: `tcp`, `udp` or `icmp`. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"address_start": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Start of the IP address range. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"address_end": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "End of the IP address range. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port_start": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Start of the port range. Forces replacement.",
			},
			"port_end": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "End of the port range. Forces replacement.",
			},
		},
	}
}

func (r *vmFirewallRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *vmFirewallRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vmFirewallRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule := client.FirewallRule{
		Action:       plan.Action.ValueString(),
		Type:         plan.Type.ValueString(),
		Comment:      plan.Comment.ValueString(),
		Protocol:     plan.Protocol.ValueString(),
		AddressStart: plan.AddressStart.ValueString(),
		AddressEnd:   plan.AddressEnd.ValueString(),
		PortStart:    plan.PortStart.ValueInt64(),
		PortEnd:      plan.PortEnd.ValueInt64(),
	}
	position, err := r.client.CreateFirewallRule(ctx, plan.ServiceID.ValueString(), rule)
	if err != nil {
		resp.Diagnostics.AddError("Error creating firewall rule", err.Error())
		return
	}

	plan.Position = types.Int64Value(position)
	plan.ID = types.StringValue(client.FormatID(plan.ServiceID.ValueString(), fmt.Sprintf("%d", position)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *vmFirewallRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vmFirewallRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, position, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	rules, err := r.client.ListFirewallRules(ctx, serviceID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading firewall rules", err.Error())
		return
	}

	var found *client.FirewallRule
	for _, rule := range rules {
		if rule.Position == position {
			r := rule
			found = &r
			break
		}
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(client.FormatID(serviceID, fmt.Sprintf("%d", position)))
	state.ServiceID = types.StringValue(serviceID)
	state.Position = types.Int64Value(position)
	state.Action = common.StringOrNull(found.Action)
	state.Type = common.StringOrNull(found.Type)
	state.Comment = common.StringOrNull(found.Comment)
	state.Protocol = common.StringOrNull(found.Protocol)
	state.AddressStart = common.StringOrNull(found.AddressStart)
	state.AddressEnd = common.StringOrNull(found.AddressEnd)
	state.PortStart = common.Int64OrNull(found.PortStart)
	state.PortEnd = common.Int64OrNull(found.PortEnd)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is unreachable: every attribute forces replacement. It exists only to
// satisfy the resource.Resource interface.
func (r *vmFirewallRuleResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"The HiTechCloud API provides no firewall rule update endpoint; the rule must be replaced instead.",
	)
}

func (r *vmFirewallRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vmFirewallRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, position, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	if err := r.client.DeleteFirewallRule(ctx, serviceID, position); client.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error deleting firewall rule", err.Error())
	}
}

func (r *vmFirewallRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := client.SplitID(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", err.Error())
		return
	}
	var position int64
	if _, err := fmt.Sscanf(parts[1], "%d", &position); err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", fmt.Sprintf("position part %q is not a number", parts[1]))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("position"), position)...)
}

func (r *vmFirewallRuleResource) parseID(_ context.Context, state *vmFirewallRuleModel, diags *diag.Diagnostics) (string, int64, bool) {
	if state.ServiceID.ValueString() != "" && !state.Position.IsNull() {
		return state.ServiceID.ValueString(), state.Position.ValueInt64(), true
	}
	parts, err := client.SplitID(state.ID.ValueString(), 2)
	if err != nil {
		diags.AddError("Unexpected Resource ID", err.Error())
		return "", 0, false
	}
	var position int64
	if _, err := fmt.Sscanf(parts[1], "%d", &position); err != nil {
		diags.AddError("Unexpected Resource ID", fmt.Sprintf("position part %q is not a number", parts[1]))
		return "", 0, false
	}
	return parts[0], position, true
}
