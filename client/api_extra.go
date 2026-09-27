// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
)

// ---------------------------------------------------------------------------
// Support tickets
// ---------------------------------------------------------------------------

// TicketDepartment is a support department.
type TicketDepartment struct {
	ID   string
	Name string
}

// Ticket is a support ticket.
type Ticket struct {
	Number      string
	Subject     string
	Status      string
	Department  string
	LastUpdated string
}

// ListTicketDepartments lists support departments (GET /api/ticket/departments).
func (c *Client) ListTicketDepartments(ctx context.Context) ([]TicketDepartment, error) {
	var raw any
	if err := c.Get(ctx, "/api/ticket/departments", nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw)
	out := make([]TicketDepartment, 0, len(list))
	for _, item := range list {
		m, ok := AsMap(item)
		if !ok {
			continue
		}
		out = append(out, TicketDepartment{
			ID:   FirstString(m, "id", "dept_id", "number"),
			Name: FirstString(m, "name", "title", "department"),
		})
	}
	return out, nil
}

// ListTickets lists support tickets (GET /api/tickets).
func (c *Client) ListTickets(ctx context.Context) ([]Ticket, error) {
	var raw any
	if err := c.Get(ctx, "/api/tickets", nil, &raw); err != nil {
		return nil, err
	}
	return parseTickets(raw), nil
}

// GetTicket fetches one ticket (GET /api/tickets/{number}).
func (c *Client) GetTicket(ctx context.Context, number string) (*Ticket, error) {
	var raw any
	if err := c.Get(ctx, "/api/tickets/"+PathEscape(number), nil, &raw); err != nil {
		return nil, err
	}
	m, ok := AsMap(unwrapValue(raw))
	if !ok {
		return nil, &APIError{StatusCode: 200, Message: "unexpected ticket payload"}
	}
	t := parseTicket(m)
	if t.Number == "" {
		t.Number = number
	}
	return &t, nil
}

// CreateTicket opens a ticket (POST /api/tickets). The ticket number is taken
// from the response or resolved by locating the newest ticket with the same
// subject.
func (c *Client) CreateTicket(ctx context.Context, deptID, subject, body string) (string, error) {
	q := Query("dept_id", deptID, "subject", subject, "body", body)
	var raw any
	if err := c.Post(ctx, "/api/tickets", q, &raw); err != nil {
		return "", err
	}
	if id := FirstString(AsMapAny(raw), "number", "ticket_number", "id"); id != "" {
		return id, nil
	}
	// Fallback: locate the ticket by subject.
	tickets, err := c.ListTickets(ctx)
	if err != nil {
		return "", err
	}
	for _, t := range tickets {
		if t.Subject == subject && t.Number != "" {
			return t.Number, nil
		}
	}
	if len(tickets) > 0 && tickets[0].Number != "" {
		return tickets[0].Number, nil
	}
	return "", &APIError{StatusCode: 200, Message: "ticket created but no identifier could be resolved"}
}

// ReplyTicket adds a reply to a ticket (POST /api/tickets/{number}).
func (c *Client) ReplyTicket(ctx context.Context, number, body string) error {
	return c.Post(ctx, "/api/tickets/"+PathEscape(number), Query("body", body), nil)
}

// CloseTicket closes a ticket (PUT /api/tickets/{number}/close).
func (c *Client) CloseTicket(ctx context.Context, number string) error {
	return c.Put(ctx, "/api/tickets/"+PathEscape(number)+"/close", nil, nil)
}

// OpenTicket re-opens a ticket (PUT /api/tickets/{number}/open).
func (c *Client) OpenTicket(ctx context.Context, number string) error {
	return c.Put(ctx, "/api/tickets/"+PathEscape(number)+"/open", nil, nil)
}

func parseTickets(raw any) []Ticket {
	list := ExtractList(raw)
	out := make([]Ticket, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseTicket(m))
		}
	}
	return out
}

func parseTicket(m map[string]any) Ticket {
	return Ticket{
		Number:      FirstString(m, "number", "ticket_number", "id"),
		Subject:     FirstString(m, "subject", "title"),
		Status:      FirstString(m, "status", "state"),
		Department:  FirstString(m, "department", "dept", "dept_name"),
		LastUpdated: FirstString(m, "last_updated", "updated_at", "date", "lastreply"),
	}
}

// ---------------------------------------------------------------------------
// Notifications & status
// ---------------------------------------------------------------------------

// Notification is a portal notification.
type Notification struct {
	ID      string
	Title   string
	Message string
	Date    string
	RelType string
	RelID   string
	Status  string
}

// ListNotifications lists portal notifications (GET /api/notifications).
func (c *Client) ListNotifications(ctx context.Context, relType, relID string) ([]Notification, error) {
	q := Query("rel_type", relType, "rel_id", relID)
	var raw any
	if err := c.Get(ctx, "/api/notifications", q, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw)
	out := make([]Notification, 0, len(list))
	for _, item := range list {
		m, ok := AsMap(item)
		if !ok {
			continue
		}
		out = append(out, Notification{
			ID:      FirstString(m, "id", "notification_id"),
			Title:   FirstString(m, "title", "subject"),
			Message: FirstString(m, "message", "body", "content"),
			Date:    FirstString(m, "date", "created_at", "time"),
			RelType: FirstString(m, "rel_type", "type"),
			RelID:   FirstString(m, "rel_id"),
			Status:  FirstString(m, "status", "state"),
		})
	}
	return out, nil
}

// StatusEntry is a service status entry (GET /api/statuses).
type StatusEntry struct {
	ID     string
	Name   string
	Status string
	Type   string
}

// ListStatuses lists service statuses (GET /api/statuses).
func (c *Client) ListStatuses(ctx context.Context) ([]StatusEntry, error) {
	var raw any
	if err := c.Get(ctx, "/api/statuses", nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw)
	out := make([]StatusEntry, 0, len(list))
	for _, item := range list {
		m, ok := AsMap(item)
		if !ok {
			continue
		}
		out = append(out, StatusEntry{
			ID:     FirstString(m, "id", "service_id"),
			Name:   FirstString(m, "name", "title", "label"),
			Status: FirstString(m, "status", "state"),
			Type:   FirstString(m, "type", "group"),
		})
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Billing extras
// ---------------------------------------------------------------------------

// Balance is the account balance (GET /api/balance).
type Balance struct {
	Credit   string
	Balance  string
	Currency string
}

// GetBalance fetches the account balance (GET /api/balance).
func (c *Client) GetBalance(ctx context.Context) (*Balance, error) {
	var raw any
	if err := c.Get(ctx, "/api/balance", nil, &raw); err != nil {
		return nil, err
	}
	m := AsMapAny(raw)
	return &Balance{
		Credit:   FirstString(m, "credit", "creditbalance"),
		Balance:  FirstString(m, "balance", "amount", "total"),
		Currency: FirstString(m, "currency"),
	}, nil
}

// PaymentFee describes the fees of a payment method.
type PaymentFee struct {
	Module string
	Name   string
	Fee    string
}

// ListPaymentFees lists payment method fees (GET /api/payment/fees).
func (c *Client) ListPaymentFees(ctx context.Context) ([]PaymentFee, error) {
	var raw any
	if err := c.Get(ctx, "/api/payment/fees", nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw)
	out := make([]PaymentFee, 0, len(list))
	for _, item := range list {
		m, ok := AsMap(item)
		if !ok {
			continue
		}
		out = append(out, PaymentFee{
			Module: FirstString(m, "module", "id", "gateway"),
			Name:   FirstString(m, "name", "title"),
			Fee:    FirstString(m, "fee", "fees", "percentage"),
		})
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// WHOIS
// ---------------------------------------------------------------------------

// GetWhois performs a WHOIS lookup for a domain (GET /api/whois/{domain}) and
// returns the flattened attributes.
func (c *Client) GetWhois(ctx context.Context, domain string) (map[string]string, error) {
	var raw any
	if err := c.Get(ctx, "/api/whois/"+PathEscape(domain), nil, &raw); err != nil {
		return nil, err
	}
	m := AsMapAny(raw)
	if flat, ok := AsMap(m["whois"]); ok {
		m = flat
	}
	return StringMap(m), nil
}

// ---------------------------------------------------------------------------
// Contacts (mutations)
// ---------------------------------------------------------------------------

// ContactInput carries the mutable fields of an account contact.
type ContactInput struct {
	Password      string
	Privileges    string
	Type          string
	CompanyName   string
	TaxID         string
	Gender        string
	LastName      string
	FirstName     string
	NationalID    string
	Email         string
	Birthday      string
	PhoneNumber   string
	Country       string
	State         string
	City          string
	Address       string
	Postcode      string
	BankName      string
	BankAccount   string
	BankAccountNm string
}

func (in ContactInput) query(includePassword bool) map[string]string {
	q := map[string]string{
		"privileges":      in.Privileges,
		"type":            in.Type,
		"companyname":     in.CompanyName,
		"taxid":           in.TaxID,
		"gender":          in.Gender,
		"lastname":        in.LastName,
		"firstname":       in.FirstName,
		"nationalid":      in.NationalID,
		"email":           in.Email,
		"birthday":        in.Birthday,
		"phonenumber":     in.PhoneNumber,
		"country":         in.Country,
		"state":           in.State,
		"city":            in.City,
		"address1":        in.Address,
		"postcode":        in.Postcode,
		"bankname":        in.BankName,
		"bankaccount":     in.BankAccount,
		"bankaccountname": in.BankAccountNm,
	}
	if includePassword {
		q["password"] = in.Password
	}
	return q
}

// CreateContact adds an account contact (POST /api/contact). The contact id is
// resolved from the response or by locating the contact by email.
func (c *Client) CreateContact(ctx context.Context, in ContactInput) (string, error) {
	var q = Query()
	for k, v := range in.query(true) {
		if v != "" {
			q.Set(k, v)
		}
	}
	var raw any
	if err := c.Post(ctx, "/api/contact", q, &raw); err != nil {
		return "", err
	}
	if id := FirstString(AsMapAny(raw), "id", "contact_id"); id != "" {
		return id, nil
	}
	contacts, err := c.ListContacts(ctx)
	if err != nil {
		return "", err
	}
	for _, ct := range contacts {
		if ct.Email == in.Email && ct.ID != "" {
			return ct.ID, nil
		}
	}
	return "", &APIError{StatusCode: 200, Message: "contact created but no identifier could be resolved"}
}

// UpdateContact edits an account contact (PUT /api/contact/{id}).
func (c *Client) UpdateContact(ctx context.Context, id string, in ContactInput) error {
	q := Query()
	for k, v := range in.query(false) {
		if v != "" {
			q.Set(k, v)
		}
	}
	return c.Put(ctx, "/api/contact/"+PathEscape(id), q, nil)
}

// ---------------------------------------------------------------------------
// Reverse DNS
// ---------------------------------------------------------------------------

// RDNSRecord maps an IP address to its reverse DNS hostname.
type RDNSRecord struct {
	IP       string
	Hostname string
}

// GetServiceRDNS lists reverse DNS entries of a service
// (GET /api/service/{id}/rdns).
func (c *Client) GetServiceRDNS(ctx context.Context, serviceID string) ([]RDNSRecord, error) {
	var raw any
	path := "/api/service/" + PathEscape(serviceID) + "/rdns"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	return parseRDNS(raw), nil
}

// parseRDNS understands both response shapes: a mapping {ip: hostname} and a
// wrapped list of records.
func parseRDNS(raw any) []RDNSRecord {
	out := []RDNSRecord{}
	switch v := unwrapValue(raw).(type) {
	case map[string]any:
		// Either a mapping {ip: hostname} or a wrapped list.
		if list := ExtractList(v); len(list) > 0 {
			for _, item := range list {
				m, ok := AsMap(item)
				if !ok {
					continue
				}
				out = append(out, RDNSRecord{
					IP:       FirstString(m, "ip", "ipaddress", "address"),
					Hostname: FirstString(m, "hostname", "rdns", "ptr", "name"),
				})
			}
			return out
		}
		for k, val := range v {
			if s := asString(val); s != "" {
				out = append(out, RDNSRecord{IP: k, Hostname: s})
			}
		}
	case []any:
		for _, item := range v {
			m, ok := AsMap(item)
			if !ok {
				continue
			}
			out = append(out, RDNSRecord{
				IP:       FirstString(m, "ip", "ipaddress", "address"),
				Hostname: FirstString(m, "hostname", "rdns", "ptr", "name"),
			})
		}
	}
	return out
}

// SetServiceRDNS sets the reverse DNS hostname of one IP address
// (POST /api/service/{id}/rdns). The API takes the IP as the query key and the
// hostname as its value; an empty hostname clears the entry.
func (c *Client) SetServiceRDNS(ctx context.Context, serviceID, ip, hostname string) error {
	path := "/api/service/" + PathEscape(serviceID) + "/rdns"
	q := Query(ip, hostname)
	return c.Post(ctx, path, q, nil)
}

// AsMapAny unwraps a payload envelope and returns it as a map (possibly
// empty), so callers can chain First* fallbacks without checking.
func AsMapAny(v any) map[string]any {
	if m, ok := AsMap(unwrapValue(v)); ok {
		return m
	}
	return map[string]any{}
}
