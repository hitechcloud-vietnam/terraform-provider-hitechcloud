// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vaifactory

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/client"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/common"
)

var (
	_ resource.Resource                = (*aiSSHKeyResource)(nil)
	_ resource.ResourceWithConfigure   = (*aiSSHKeyResource)(nil)
	_ resource.ResourceWithImportState = (*aiSSHKeyResource)(nil)
)

// AISSHKeyResource returns the hitechcloud_ai_ssh_key resource.
func AISSHKeyResource() resource.Resource {
	return &aiSSHKeyResource{}
}

type aiSSHKeyResource struct {
	client *client.Client
}

type aiSSHKeyModel struct {
	ID        types.String `tfsdk:"id"`
	ServiceID types.String `tfsdk:"service_id"`
	KeyID     types.String `tfsdk:"key_id"`
	Name      types.String `tfsdk:"name"`
	PublicKey types.String `tfsdk:"public_key"`
	Default   types.Bool   `tfsdk:"default"`
	IsDefault types.Bool   `tfsdk:"is_default"`
}

func (r *aiSSHKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ai_ssh_key"
}

func (r *aiSSHKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpenSSH public key registered with the HiTechCloud " +
			"AI Factory service (`POST /api/service/{service_id}/sshkeys`). The API has " +
			"no key update endpoint, so name and key content force replacement. " +
			"Setting `default = true` marks the key as the service default.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `<service_id>/<key_id>`.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"service_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "HostBill service id of the AI Factory service. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"key_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SSH key identifier.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Key name. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"public_key": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "OpenSSH public key material. Forces replacement.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"default": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Set to `true` to make this the default SSH key of the service (applied on create/update).",
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"is_default": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the key is currently the service default (as reported by the API).",
			},
		},
	}
}

func (r *aiSSHKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *aiSSHKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan aiSSHKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	keyID, err := r.client.CreateAISSHKey(ctx, plan.ServiceID.ValueString(), plan.Name.ValueString(), plan.PublicKey.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating SSH key", err.Error())
		return
	}
	plan.KeyID = types.StringValue(keyID)
	plan.ID = types.StringValue(client.FormatID(plan.ServiceID.ValueString(), keyID))

	if !plan.Default.IsNull() && !plan.Default.IsUnknown() && plan.Default.ValueBool() {
		if err := r.client.SetDefaultAISSHKey(ctx, plan.ServiceID.ValueString(), keyID); err != nil {
			resp.Diagnostics.AddError("Error setting default SSH key", err.Error())
			return
		}
	}

	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *aiSSHKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state aiSSHKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, keyID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}
	state.ServiceID = types.StringValue(serviceID)
	state.KeyID = types.StringValue(keyID)
	if !r.refresh(ctx, &state, &resp.Diagnostics) {
		resp.State.RemoveResource(ctx)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update only reacts to `default = true`; everything else forces replacement.
func (r *aiSSHKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan aiSSHKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, keyID, ok := r.parseID(ctx, &plan, &resp.Diagnostics)
	if !ok {
		return
	}

	if !plan.Default.IsNull() && !plan.Default.IsUnknown() && plan.Default.ValueBool() {
		if err := r.client.SetDefaultAISSHKey(ctx, serviceID, keyID); err != nil {
			resp.Diagnostics.AddError("Error setting default SSH key", err.Error())
			return
		}
	}

	plan.ID = types.StringValue(client.FormatID(serviceID, keyID))
	r.refresh(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *aiSSHKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state aiSSHKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceID, keyID, ok := r.parseID(ctx, &state, &resp.Diagnostics)
	if !ok {
		return
	}

	if err := r.client.DeleteAISSHKey(ctx, serviceID, keyID); client.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error deleting SSH key", err.Error())
	}
}

func (r *aiSSHKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts, err := client.SplitID(req.ID, 2)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Import Identifier", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key_id"), parts[1])...)
}

func (r *aiSSHKeyResource) parseID(_ context.Context, state *aiSSHKeyModel, diags *diag.Diagnostics) (string, string, bool) {
	if state.ServiceID.ValueString() != "" && state.KeyID.ValueString() != "" {
		return state.ServiceID.ValueString(), state.KeyID.ValueString(), true
	}
	parts, err := client.SplitID(state.ID.ValueString(), 2)
	if err != nil {
		diags.AddError("Unexpected Resource ID", err.Error())
		return "", "", false
	}
	return parts[0], parts[1], true
}

// refresh reloads the key; returns false when it no longer exists. The
// `default` attribute is intentionally not refreshed (it tracks user intent).
func (r *aiSSHKeyResource) refresh(ctx context.Context, model *aiSSHKeyModel, diags *diag.Diagnostics) bool {
	key, err := r.client.GetAISSHKey(ctx, model.ServiceID.ValueString(), model.KeyID.ValueString())
	if client.IsNotFound(err) {
		return false
	}
	if err != nil {
		diags.AddError("Error reading SSH key", err.Error())
		return true
	}
	model.Name = common.StringOrNull(key.Name)
	model.PublicKey = common.StringOrNull(key.PublicKey)
	model.IsDefault = types.BoolValue(key.IsDefault)
	return true
}
