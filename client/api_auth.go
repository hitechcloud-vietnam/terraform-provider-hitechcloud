// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
)

// LoginResult is the outcome of POST /api/login.
type LoginResult struct {
	Token        string
	RefreshToken string
}

// Login performs POST /api/login and stores the returned bearer token on the
// client. username and password travel as query parameters, per the API
// specification.
func (c *Client) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	var raw any
	if err := c.Post(ctx, "/api/login", Query("username", username, "password", password), &raw); err != nil {
		return nil, err
	}

	result := &LoginResult{
		Token: FirstString(raw, "token", "access_token", "auth_token", "api_token", "access"),
		RefreshToken: FirstString(raw,
			"refresh_token", "refreshToken", "refresh", "refresh_access_token"),
	}
	if result.Token == "" {
		// Some deployments nest credentials under a user/account object.
		if nested := First(raw, "user", "account", "client"); nested != nil {
			result.Token = FirstString(nested, "token", "access_token", "auth_token")
		}
	}
	if result.Token == "" {
		return nil, fmt.Errorf("login response did not include a token")
	}
	c.SetToken(result.Token)
	if result.RefreshToken != "" {
		c.SetRefreshToken(result.RefreshToken)
	}
	return result, nil
}

// Logout performs POST /api/logout, invalidating the current token.
func (c *Client) Logout(ctx context.Context) error {
	return c.Post(ctx, "/api/logout", nil, nil)
}

// RefreshToken exchanges a refresh token for a new access token using
// POST /api/token.
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*LoginResult, error) {
	var raw any
	if err := c.Post(ctx, "/api/token", Query("refresh_token", refreshToken), &raw); err != nil {
		return nil, err
	}
	result := &LoginResult{
		Token:        FirstString(raw, "token", "access_token", "auth_token", "api_token"),
		RefreshToken: FirstString(raw, "refresh_token", "refreshToken", "refresh"),
	}
	if result.Token == "" {
		return nil, fmt.Errorf("token refresh response did not include a token")
	}
	c.SetToken(result.Token)
	if result.RefreshToken != "" {
		c.SetRefreshToken(result.RefreshToken)
	}
	return result, nil
}

// RevokeToken invalidates an authorization/refresh token via POST /api/revoke.
func (c *Client) RevokeToken(ctx context.Context, refreshToken string) error {
	return c.Post(ctx, "/api/revoke", Query("refresh_token", refreshToken), nil)
}

// Account mirrors the registration details returned by GET /api/details.
type Account struct {
	ID            string
	Email         string
	Type          string
	FirstName     string
	LastName      string
	CompanyName   string
	TaxID         string
	Gender        string
	NationalID    string
	Birthday      string
	PhoneNumber   string
	Country       string
	State         string
	City          string
	Address       string
	Postcode      string
	Currency      string
	BankName      string
	BankAccount   string
	BankAccountNm string
}

// GetAccount returns registration details for the authenticated account
// (GET /api/details).
func (c *Client) GetAccount(ctx context.Context) (*Account, error) {
	var raw any
	if err := c.Get(ctx, "/api/details", nil, &raw); err != nil {
		return nil, err
	}
	return parseAccount(raw), nil
}

func parseAccount(raw any) *Account {
	return &Account{
		ID:            FirstString(raw, "id", "client_id", "user_id", "account_id"),
		Email:         FirstString(raw, "email", "email_address", "mail"),
		Type:          FirstString(raw, "type", "account_type", "client_type"),
		FirstName:     FirstString(raw, "firstname", "first_name", "fname"),
		LastName:      FirstString(raw, "lastname", "last_name", "lname"),
		CompanyName:   FirstString(raw, "companyname", "company_name", "company", "organization"),
		TaxID:         FirstString(raw, "taxid", "tax_id", "vat"),
		Gender:        FirstString(raw, "gender"),
		NationalID:    FirstString(raw, "nationalid", "national_id"),
		Birthday:      FirstString(raw, "birthday", "birth_date", "dob"),
		PhoneNumber:   FirstString(raw, "phonenumber", "phone_number", "phone"),
		Country:       FirstString(raw, "country"),
		State:         FirstString(raw, "state"),
		City:          FirstString(raw, "city"),
		Address:       FirstString(raw, "address1", "address", "address_1"),
		Postcode:      FirstString(raw, "postcode", "zip", "postal_code"),
		Currency:      FirstString(raw, "currency"),
		BankName:      FirstString(raw, "bankname", "bank_name"),
		BankAccount:   FirstString(raw, "bankaccount", "bank_account"),
		BankAccountNm: FirstString(raw, "bankaccountname", "bank_account_name"),
	}
}
