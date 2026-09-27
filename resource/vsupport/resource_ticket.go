// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vsupport

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
	_ resource.Resource                = (*ticketResource)(nil)
	_ resource.ResourceWithConfigure   = (*ticketResource)(nil)
	_ resource.ResourceWithImportState = (*ticketResource)(nil)
)

// TicketResource returns the hitechcloud_ticket resource.
func TicketResource() resource.Resource {
	return &ticketResource{}
}

type ticketResource struct {
	client *client.Client
}

type ticketModel struct {
	ID          types.String `tfsdk:"id"`
	Number      types.String `tfsdk:"number"`
	DeptID      types.String `tfsdk:"dept_id"`
	Subject     types.String `tfsdk:"subject"`
	Body        types.String `tfsdk:"body"`
	Status      types.String `tfsdk:"status"`
	Department  types.String `tfsdk:"department"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

func (r *ticketResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ticket"
}

func (r *ticketResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a HiTechCloud support ticket (`POST /api/tickets`). " +
			"Destroying the resource closes the ticket. Replies are not managed by " +
			"Terraform; the `body` is only sent when the ticket is opened.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Ticket number.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"number": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Ticket number as returned by the API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dept_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Support department identifier. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subject": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Ticket subject. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"body": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Initial ticket body. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Ticket status.",
			},
			"department": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Support department name.",
			},
			"last_updated": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last update timestamp.",
			},
		},
	}
}

func (r *ticketResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ticketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ticketModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	number, err := r.client.CreateTicket(ctx,
		plan.DeptID.ValueString(),
		plan.Subject.ValueString(),
		plan.Body.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error creating ticket", err.Error())
		return
	}

	plan.ID = types.StringValue(number)
	plan.Number = types.StringValue(number)
	if !r.refresh(ctx, &plan) {
		resp.Diagnostics.AddError("Error reading ticket", fmt.Sprintf("ticket %q not found after creation", number))
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ticketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ticketModel
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

func (r *ticketResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Ticket Update Not Supported",
		"The API does not allow editing ticket fields; all attributes force a new resource.",
	)
}

func (r *ticketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ticketModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The API has no ticket deletion; closing is the closest equivalent.
	if err := r.client.CloseTicket(ctx, state.Number.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error closing ticket", err.Error())
	}
}

func (r *ticketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// refresh reloads the ticket into the model and reports whether it exists.
func (r *ticketResource) refresh(ctx context.Context, model *ticketModel) bool {
	t, err := r.client.GetTicket(ctx, model.Number.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			return false
		}
		return true // keep state; transient errors surface on next plan
	}
	model.ID = types.StringValue(t.Number)
	model.Number = types.StringValue(t.Number)
	model.Status = common.StringOrNull(t.Status)
	model.Department = common.StringOrNull(t.Department)
	model.LastUpdated = common.StringOrNull(t.LastUpdated)
	return true
}
