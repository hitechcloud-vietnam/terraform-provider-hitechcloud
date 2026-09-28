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
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/common"
)

var (
	_ resource.Resource                = (*dnssecKeyResource)(nil)
	_ resource.ResourceWithConfigure   = (*dnssecKeyResource)(nil)
	_ resource.ResourceWithImportState = (*dnssecKeyResource)(nil)
)

// DNSSECKeyResource returns the hitechcloud_dnssec_key resource.
func DNSSECKeyResource() resource.Resource {
	return &dnssecKeyResource{}
}

type dnssecKeyResource struct {
	client *client.Client
}

type dnssecKeyModel struct {
	ID         types.String `tfsdk:"id"`
	DomainID   types.String `tfsdk:"domain_id"`
	KeyTag     types.String `tfsdk:"key_tag"`
	Algorithm  types.String `tfsdk:"algorithm"`
	DigestType types.String `tfsdk:"digest_type"`
	Digest     types.String `tfsdk:"digest"`
	Flags      types.String `tfsdk:"flags"`
	Protocol   types.String `tfsdk:"protocol"`
	PublicKey  types.String `tfsdk:"public_key"`
}

func (r *dnssecKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dnssec_key"
}

func (r *dnssecKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a DNSSEC key (DS record) of a registered domain " +
			"(`PUT /api/domain/{id}/dnssec`). Destroying the resource removes the key " +
			"with `DELETE /api/domain/{id}/dnssec/{key}`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier `domain_id/key_tag`.",
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
			"key_tag": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "DS key tag. Changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"algorithm": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "DNSSEC algorithm number, for example `8` (RSA/SHA-256) or `13` (ECDSA P-256).",
			},
			"digest_type": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "DS digest type, for example `2` (SHA-256).",
			},
			"digest": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Hex-encoded DS digest of the DNSKEY record.",
			},
			"flags": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "DNSKEY flags, usually `257` (KSK) or `256` (ZSK).",
			},
			"protocol": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "DNSKEY protocol number, usually `3`.",
			},
			"public_key": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Base64-encoded DNSKEY public key material.",
			},
		},
	}
}

func (r *dnssecKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// addParams builds the query parameters accepted by the DNSSEC add endpoint.
func (m *dnssecKeyModel) addParams() map[string]string {
	params := map[string]string{
		"key_tag":   m.KeyTag.ValueString(),
		"algorithm": m.Algorithm.ValueString(),
	}
	if v := m.DigestType.ValueString(); v != "" {
		params["digest_type"] = v
	}
	if v := m.Digest.ValueString(); v != "" {
		params["digest"] = v
	}
	if v := m.Flags.ValueString(); v != "" {
		params["flags"] = v
	}
	if v := m.Protocol.ValueString(); v != "" {
		params["protocol"] = v
	}
	if v := m.PublicKey.ValueString(); v != "" {
		params["public_key"] = v
	}
	return params
}

func (r *dnssecKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan dnssecKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.AddDNSSECKey(ctx, plan.DomainID.ValueString(), plan.addParams()); err != nil {
		resp.Diagnostics.AddError("Error adding DNSSEC key", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.DomainID.ValueString(), plan.KeyTag.ValueString()))
	r.refresh(ctx, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dnssecKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state dnssecKeyModel
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

func (r *dnssecKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan dnssecKeyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The API has no in-place key edit: replace the key under the same tag.
	if err := r.client.RemoveDNSSECKey(ctx, plan.DomainID.ValueString(), plan.KeyTag.ValueString()); err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error removing old DNSSEC key", err.Error())
		return
	}
	if err := r.client.AddDNSSECKey(ctx, plan.DomainID.ValueString(), plan.addParams()); err != nil {
		resp.Diagnostics.AddError("Error updating DNSSEC key", err.Error())
		return
	}

	r.refresh(ctx, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *dnssecKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state dnssecKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.RemoveDNSSECKey(ctx, state.DomainID.ValueString(), state.KeyTag.ValueString()); err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error removing DNSSEC key", err.Error())
	}
}

func (r *dnssecKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	parts, err := client.SplitID(id, 2)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf(
			"Expected import ID in the format domain_id/key_tag: %s", err))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key_tag"), parts[1])...)
}

// refresh reloads the DNSSEC key list and reports whether the key still exists.
func (r *dnssecKeyResource) refresh(ctx context.Context, model *dnssecKeyModel) bool {
	keys, err := r.client.ListDNSSECKeys(ctx, model.DomainID.ValueString())
	if err != nil {
		return !client.IsNotFound(err)
	}
	for _, item := range keys {
		m, ok := client.AsMap(item)
		if !ok {
			continue
		}
		if client.FirstString(m, "key_tag", "keytag", "tag") != model.KeyTag.ValueString() {
			continue
		}
		model.Algorithm = common.StringOrNull(client.FirstString(m, "algorithm", "alg"))
		model.DigestType = common.StringOrNull(client.FirstString(m, "digest_type", "digesttype"))
		model.Digest = common.StringOrNull(client.FirstString(m, "digest"))
		model.Flags = common.StringOrNull(client.FirstString(m, "flags", "flag"))
		model.Protocol = common.StringOrNull(client.FirstString(m, "protocol", "proto"))
		model.PublicKey = common.StringOrNull(client.FirstString(m, "public_key", "pubkey", "key"))
		return true
	}
	return false
}
