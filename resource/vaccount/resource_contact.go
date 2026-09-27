// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package vaccount

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
	_ resource.Resource                = (*contactResource)(nil)
	_ resource.ResourceWithConfigure   = (*contactResource)(nil)
	_ resource.ResourceWithImportState = (*contactResource)(nil)
)

// ContactResource returns the hitechcloud_contact resource.
func ContactResource() resource.Resource {
	return &contactResource{}
}

type contactResource struct {
	client *client.Client
}

type contactModel struct {
	ID            types.String `tfsdk:"id"`
	ContactID     types.String `tfsdk:"contact_id"`
	Email         types.String `tfsdk:"email"`
	Password      types.String `tfsdk:"password"`
	Privileges    types.String `tfsdk:"privileges"`
	Type          types.String `tfsdk:"type"`
	FirstName     types.String `tfsdk:"first_name"`
	LastName      types.String `tfsdk:"last_name"`
	CompanyName   types.String `tfsdk:"company_name"`
	TaxID         types.String `tfsdk:"tax_id"`
	Gender        types.String `tfsdk:"gender"`
	NationalID    types.String `tfsdk:"national_id"`
	Birthday      types.String `tfsdk:"birthday"`
	PhoneNumber   types.String `tfsdk:"phone_number"`
	Country       types.String `tfsdk:"country"`
	State         types.String `tfsdk:"state"`
	City          types.String `tfsdk:"city"`
	Address       types.String `tfsdk:"address"`
	Postcode      types.String `tfsdk:"postcode"`
	BankName      types.String `tfsdk:"bank_name"`
	BankAccount   types.String `tfsdk:"bank_account"`
	BankAccountNm types.String `tfsdk:"bank_account_name"`
}

func (r *contactResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contact"
}

func (r *contactResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a contact of the HiTechCloud account " +
			"(`POST/PUT /api/contact`). The API has no contact deletion, so " +
			"destroying the resource only removes it from state (the contact " +
			"remains in the account).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Composite identifier (same as `contact_id`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"contact_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Identifier of the contact as returned by the API.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"email": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Email address of the contact.",
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				MarkdownDescription: "Password for the contact login. Only used at creation; " +
					"changing this forces a new resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"privileges":   schema.StringAttribute{Optional: true, MarkdownDescription: "Privilege set of the contact."},
			"type":         schema.StringAttribute{Optional: true, MarkdownDescription: "Contact type."},
			"first_name":   schema.StringAttribute{Optional: true, MarkdownDescription: "First name."},
			"last_name":    schema.StringAttribute{Optional: true, MarkdownDescription: "Last name."},
			"company_name": schema.StringAttribute{Optional: true, MarkdownDescription: "Company name."},
			"tax_id":       schema.StringAttribute{Optional: true, MarkdownDescription: "Tax identifier."},
			"gender":       schema.StringAttribute{Optional: true, MarkdownDescription: "Gender."},
			"national_id":  schema.StringAttribute{Optional: true, MarkdownDescription: "National ID."},
			"birthday":     schema.StringAttribute{Optional: true, MarkdownDescription: "Birthday."},
			"phone_number": schema.StringAttribute{Optional: true, MarkdownDescription: "Phone number."},
			"country":      schema.StringAttribute{Optional: true, MarkdownDescription: "Country."},
			"state":        schema.StringAttribute{Optional: true, MarkdownDescription: "State / province."},
			"city":         schema.StringAttribute{Optional: true, MarkdownDescription: "City."},
			"address":      schema.StringAttribute{Optional: true, MarkdownDescription: "Street address."},
			"postcode":     schema.StringAttribute{Optional: true, MarkdownDescription: "Postal code."},
			"bank_name":    schema.StringAttribute{Optional: true, MarkdownDescription: "Bank name."},
			"bank_account": schema.StringAttribute{Optional: true, MarkdownDescription: "Bank account number."},
			"bank_account_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Bank account holder name.",
			},
		},
	}
}

func (r *contactResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (m *contactModel) input() client.ContactInput {
	return client.ContactInput{
		Password:      m.Password.ValueString(),
		Privileges:    m.Privileges.ValueString(),
		Type:          m.Type.ValueString(),
		CompanyName:   m.CompanyName.ValueString(),
		TaxID:         m.TaxID.ValueString(),
		Gender:        m.Gender.ValueString(),
		LastName:      m.LastName.ValueString(),
		FirstName:     m.FirstName.ValueString(),
		NationalID:    m.NationalID.ValueString(),
		Email:         m.Email.ValueString(),
		Birthday:      m.Birthday.ValueString(),
		PhoneNumber:   m.PhoneNumber.ValueString(),
		Country:       m.Country.ValueString(),
		State:         m.State.ValueString(),
		City:          m.City.ValueString(),
		Address:       m.Address.ValueString(),
		Postcode:      m.Postcode.ValueString(),
		BankName:      m.BankName.ValueString(),
		BankAccount:   m.BankAccount.ValueString(),
		BankAccountNm: m.BankAccountNm.ValueString(),
	}
}

func (r *contactResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan contactModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateContact(ctx, plan.input())
	if err != nil {
		resp.Diagnostics.AddError("Error creating contact", err.Error())
		return
	}

	plan.ID = types.StringValue(id)
	plan.ContactID = types.StringValue(id)
	r.refresh(ctx, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *contactResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state contactModel
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

func (r *contactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan contactModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateContact(ctx, plan.ContactID.ValueString(), plan.input()); err != nil {
		resp.Diagnostics.AddError("Error updating contact", err.Error())
		return
	}

	r.refresh(ctx, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *contactResource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// The HiTechCloud API has no contact deletion endpoint. Removing the
	// resource from state leaves the contact in the account.
	resp.Diagnostics.AddWarning(
		"Contact Not Deleted",
		"The HiTechCloud API does not support deleting contacts; the contact remains in the account.",
	)
}

func (r *contactResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// refresh reloads the contact and reports whether it still exists.
func (r *contactResource) refresh(ctx context.Context, model *contactModel) bool {
	ct, err := r.client.GetContact(ctx, model.ContactID.ValueString())
	if err != nil {
		return !client.IsNotFound(err)
	}
	model.ID = types.StringValue(ct.ID)
	model.ContactID = types.StringValue(ct.ID)
	model.Email = common.StringOrNull(ct.Email)
	model.FirstName = common.StringOrNull(ct.FirstName)
	model.LastName = common.StringOrNull(ct.LastName)
	model.CompanyName = common.StringOrNull(ct.CompanyName)
	model.PhoneNumber = common.StringOrNull(ct.PhoneNumber)
	return true
}
