// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vportal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/client"
	"github.com/hitechcloud-vietnam/terraform-provider-hitechcloud/resource/common"
)

var (
	_ resource.Resource              = (*urlShortenerLinkResource)(nil)
	_ resource.ResourceWithConfigure = (*urlShortenerLinkResource)(nil)
)

// URLShortenerLinkResource returns the hitechcloud_url_shortener_link resource.
func URLShortenerLinkResource() resource.Resource {
	return &urlShortenerLinkResource{}
}

type urlShortenerLinkResource struct {
	client *client.Client
}

type urlShortenerLinkModel struct {
	ID       types.String `tfsdk:"id"`
	URL      types.String `tfsdk:"url"`
	Label    types.String `tfsdk:"label"`
	ShortURL types.String `tfsdk:"short_url"`
	Clicks   types.Int64  `tfsdk:"clicks"`
}

func (r *urlShortenerLinkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_url_shortener_link"
}

func (r *urlShortenerLinkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a HiTechCloud short link " +
			"(`POST /api/url-shortener/shorten`, `GET/DELETE /api/url-shortener/links/{id}`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Short link identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Long URL to shorten. Forces replacement (the API has no link update endpoint).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"label": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Optional label for the short link. Forces replacement.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"short_url": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The generated short URL.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"clicks": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Number of times the short link has been visited.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *urlShortenerLinkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *urlShortenerLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan urlShortenerLinkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := r.client.ShortenURL(ctx, plan.URL.ValueString(), plan.Label.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating short link", err.Error())
		return
	}

	plan.ID = types.StringValue(link.ID)
	plan.URL = types.StringValue(link.URL)
	plan.Label = common.StringOrNull(link.Label)
	plan.ShortURL = common.StringOrNull(link.Short)
	plan.Clicks = types.Int64Value(link.Clicks)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *urlShortenerLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state urlShortenerLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	link, err := r.client.GetURLLink(ctx, state.ID.ValueString())
	if client.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading short link", err.Error())
		return
	}

	state.URL = common.StringOrNull(link.URL)
	state.Label = common.StringOrNull(link.Label)
	state.ShortURL = common.StringOrNull(link.Short)
	state.Clicks = types.Int64Value(link.Clicks)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update is unreachable: every attribute forces replacement. It exists only to
// satisfy the resource.Resource interface.
func (r *urlShortenerLinkResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"Short links cannot be edited through the HiTechCloud API; the resource must be replaced instead.",
	)
}

func (r *urlShortenerLinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state urlShortenerLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteURLLink(ctx, state.ID.ValueString()); client.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error deleting short link", err.Error())
	}
}
