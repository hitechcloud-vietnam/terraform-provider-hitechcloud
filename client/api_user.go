// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
)

// User Profile (folder "User Profile") and account-level helpers.

// UpdateAccount updates the registration details of the authenticated
// account (PUT /api/details). params accepts: type, companyname, taxid,
// gender, lastname, firstname, nationalid, email, birthday, phonenumber,
// country, state, city, address1, postcode, bankname, bankaccount,
// bankaccountname.
func (c *Client) UpdateAccount(ctx context.Context, params map[string]string) error {
	_, err := c.putMap(ctx, "/api/details", qmap(params))
	return err
}

// GetAccountLogs returns the account activity log (GET /api/logs).
func (c *Client) GetAccountLogs(ctx context.Context) ([]any, error) {
	return c.getList(ctx, "/api/logs", nil)
}

// PasswordReset requests a password reset email (POST /api/passwordreset).
func (c *Client) PasswordReset(ctx context.Context, email string) error {
	_, err := c.postMap(ctx, "/api/passwordreset", Query("email", email))
	return err
}

// Signup registers a new account (POST /api/signup). params accepts: type,
// companyname, taxid, gender, lastname, firstname, nationalid, email,
// birthday, phonenumber, country, state, city, address1, password, postcode,
// currency. The created account (and its credentials) is returned.
func (c *Client) Signup(ctx context.Context, params map[string]string) (map[string]any, error) {
	return c.postMap(ctx, "/api/signup", qmap(params))
}

// GetContactPrivileges lists the privileges that can be granted to contacts
// (GET /api/contact/privileges).
func (c *Client) GetContactPrivileges(ctx context.Context) (map[string]any, error) {
	return c.getMap(ctx, "/api/contact/privileges", nil)
}
