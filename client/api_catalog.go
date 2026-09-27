// Copyright (c) HiTechCloud Vietnam. SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/url"
)

// Service is an entry of GET /api/service.
type Service struct {
	ID           string
	Name         string
	Status       string
	Label        string
	Domain       string
	Group        string
	Cycle        string
	Amount       string
	NextDue      string
	IP           string
	RegisterDate string
}

// ListServices returns all services under the account (GET /api/service).
func (c *Client) ListServices(ctx context.Context) ([]Service, error) {
	var raw any
	if err := c.Get(ctx, "/api/service", nil, &raw); err != nil {
		return nil, err
	}
	return parseServices(ExtractList(raw, "services", "accounts")), nil
}

// GetService returns details for a single service (GET /api/service/:id).
func (c *Client) GetService(ctx context.Context, id string) (*Service, error) {
	var raw any
	if err := c.Get(ctx, "/api/service/"+PathEscape(id), nil, &raw); err != nil {
		return nil, err
	}
	m, ok := AsMap(raw)
	if !ok {
		if list := ExtractList(raw, "services"); len(list) > 0 {
			m, _ = AsMap(list[0])
		}
	}
	if m == nil {
		return nil, fmt.Errorf("unexpected response shape for service %s", id)
	}
	svc := parseService(m)
	if svc.ID == "" {
		svc.ID = id
	}
	return &svc, nil
}

func parseServices(list []any) []Service {
	out := make([]Service, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseService(m))
		}
	}
	return out
}

func parseService(m map[string]any) Service {
	return Service{
		ID:           FirstString(m, "id", "service_id", "account_id"),
		Name:         FirstString(m, "name", "product", "product_name", "title"),
		Status:       FirstString(m, "status", "state"),
		Label:        FirstString(m, "label", "custom_label"),
		Domain:       FirstString(m, "domain", "domain_name", "hostname"),
		Group:        FirstString(m, "group", "product_group", "category", "category_name"),
		Cycle:        FirstString(m, "cycle", "billing_cycle", "period"),
		Amount:       FirstString(m, "amount", "price", "total"),
		NextDue:      FirstString(m, "next_due", "next_due_date", "duedate", "due_date"),
		IP:           FirstString(m, "ip", "dedicated_ip", "ip_address"),
		RegisterDate: FirstString(m, "registerdate", "register_date", "created", "date_created"),
	}
}

// Domain is an entry of GET /api/domain.
type Domain struct {
	ID            string
	Name          string
	Status        string
	Registrar     string
	RegDate       string
	Expires       string
	Autorenew     bool
	RegistrarLock bool
	IDProtection  bool
	Nameservers   []string
	EPPCode       string
}

// ListDomains returns all domains under the account (GET /api/domain).
func (c *Client) ListDomains(ctx context.Context) ([]Domain, error) {
	var raw any
	if err := c.Get(ctx, "/api/domain", nil, &raw); err != nil {
		return nil, err
	}
	return parseDomains(ExtractList(raw, "domains")), nil
}

// GetDomain returns domain details by numeric id (GET /api/domain/:id).
func (c *Client) GetDomain(ctx context.Context, id string) (*Domain, error) {
	var raw any
	if err := c.Get(ctx, "/api/domain/"+PathEscape(id), nil, &raw); err != nil {
		return nil, err
	}
	return parseDomainResponse(raw, id), nil
}

// GetDomainByName returns domain details by domain name
// (GET /api/domain/name/:name).
func (c *Client) GetDomainByName(ctx context.Context, name string) (*Domain, error) {
	var raw any
	if err := c.Get(ctx, "/api/domain/name/"+PathEscape(name), nil, &raw); err != nil {
		return nil, err
	}
	return parseDomainResponse(raw, ""), nil
}

func parseDomainResponse(raw any, fallbackID string) *Domain {
	m, ok := AsMap(raw)
	if !ok {
		if list := ExtractList(raw, "domains"); len(list) > 0 {
			m, _ = AsMap(list[0])
		}
	}
	if m == nil {
		return &Domain{ID: fallbackID}
	}
	d := parseDomain(m)
	if d.ID == "" {
		d.ID = fallbackID
	}
	return &d
}

func parseDomains(list []any) []Domain {
	out := make([]Domain, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseDomain(m))
		}
	}
	return out
}

func parseDomain(m map[string]any) Domain {
	return Domain{
		ID:            FirstString(m, "id", "domain_id"),
		Name:          FirstString(m, "name", "domain", "domain_name"),
		Status:        FirstString(m, "status", "state"),
		Registrar:     FirstString(m, "registrar", "registrar_name"),
		RegDate:       FirstString(m, "regdate", "registration_date", "registered", "created"),
		Expires:       FirstString(m, "expires", "expiry", "expiry_date", "expiration_date", "expires_date"),
		Autorenew:     FirstBool(m, "autorenew", "auto_renew", "autorenew_enabled"),
		RegistrarLock: FirstBool(m, "reglock", "registrar_lock", "transfer_lock", "locked"),
		IDProtection:  FirstBool(m, "idprotection", "id_protection", "idprotect", "privacy"),
		Nameservers:   StringList(First(m, "nameservers", "ns", "nameserver", "nameservers_list")),
		EPPCode:       FirstString(m, "epp", "epp_code", "transfer_secret"),
	}
}

// TLD describes an entry of GET /api/domain/order (available TLDs).
type TLD struct {
	ID       string
	TLD      string
	MinYears int64
	MaxYears int64
}

// ListDomainTLDs returns TLDs available for registration and transfer
// (GET /api/domain/order).
func (c *Client) ListDomainTLDs(ctx context.Context) ([]TLD, error) {
	var raw any
	if err := c.Get(ctx, "/api/domain/order", nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "tlds", "tld", "domains")
	out := make([]TLD, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, TLD{
				ID:       FirstString(m, "id", "tld_id"),
				TLD:      FirstString(m, "tld", "name", "extension"),
				MinYears: FirstInt64(m, "min_years", "minyears", "min"),
				MaxYears: FirstInt64(m, "max_years", "maxyears", "max"),
			})
		}
	}
	return out, nil
}

// CheckDomainAvailability reports whether a domain can be registered
// (POST /api/domain/lookup). The API answers with status "ok" when the domain
// is available and an empty response otherwise.
func (c *Client) CheckDomainAvailability(ctx context.Context, name string) (bool, string, error) {
	var raw any
	if err := c.Post(ctx, "/api/domain/lookup", Query("name", name), &raw); err != nil {
		return false, "", err
	}
	status := FirstString(raw, "status", "result", "availability")
	if status == "" {
		if b, ok := asBool(raw); ok {
			return b, "", nil
		}
	}
	return status == "ok" || status == "available" || status == "true", status, nil
}

// Invoice is an entry of GET /api/invoice.
type Invoice struct {
	ID       string
	Number   string
	Status   string
	Currency string
	Total    string
	Date     string
	DueDate  string
}

// ListInvoices returns all invoices under the account (GET /api/invoice).
func (c *Client) ListInvoices(ctx context.Context) ([]Invoice, error) {
	var raw any
	if err := c.Get(ctx, "/api/invoice", nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "invoices")
	out := make([]Invoice, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseInvoice(m))
		}
	}
	return out, nil
}

// GetInvoice returns invoice details (GET /api/invoice/:id).
func (c *Client) GetInvoice(ctx context.Context, id string) (*Invoice, error) {
	var raw any
	if err := c.Get(ctx, "/api/invoice/"+PathEscape(id), nil, &raw); err != nil {
		return nil, err
	}
	m, ok := AsMap(raw)
	if !ok {
		return nil, fmt.Errorf("unexpected response shape for invoice %s", id)
	}
	inv := parseInvoice(m)
	if inv.ID == "" {
		inv.ID = id
	}
	return &inv, nil
}

func parseInvoice(m map[string]any) Invoice {
	return Invoice{
		ID:       FirstString(m, "id", "invoice_id"),
		Number:   FirstString(m, "number", "invoice_number", "num"),
		Status:   FirstString(m, "status", "state"),
		Currency: FirstString(m, "currency"),
		Total:    FirstString(m, "total", "total_due", "amount", "balance"),
		Date:     FirstString(m, "date", "date_created", "created"),
		DueDate:  FirstString(m, "duedate", "due_date", "date_due"),
	}
}

// PaymentMethod is an entry of GET /api/payment.
type PaymentMethod struct {
	ID   string
	Name string
}

// ListPaymentMethods returns available payment methods (GET /api/payment).
func (c *Client) ListPaymentMethods(ctx context.Context) ([]PaymentMethod, error) {
	var raw any
	if err := c.Get(ctx, "/api/payment", nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "payments", "payment_methods", "methods", "gateways")
	out := make([]PaymentMethod, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, PaymentMethod{
				ID:   FirstString(m, "id", "module", "gateway", "key"),
				Name: FirstString(m, "name", "display_name", "title"),
			})
		}
	}
	return out, nil
}

// ProductCategory is an entry of GET /api/category.
type ProductCategory struct {
	ID          string
	Name        string
	Description string
}

// ListProductCategories returns product categories (GET /api/category).
func (c *Client) ListProductCategories(ctx context.Context) ([]ProductCategory, error) {
	var raw any
	if err := c.Get(ctx, "/api/category", nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "categories")
	out := make([]ProductCategory, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, ProductCategory{
				ID:          FirstString(m, "id", "category_id"),
				Name:        FirstString(m, "name", "title"),
				Description: FirstString(m, "description", "desc"),
			})
		}
	}
	return out, nil
}

// Product is an entry of GET /api/category/:category_id/product.
type Product struct {
	ID          string
	Name        string
	Description string
	Cycle       string
	Price       string
}

// ListProducts returns purchasable products in a category
// (GET /api/category/:category_id/product).
func (c *Client) ListProducts(ctx context.Context, categoryID string) ([]Product, error) {
	var raw any
	path := "/api/category/" + PathEscape(categoryID) + "/product"
	if err := c.Get(ctx, path, nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "products")
	out := make([]Product, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, Product{
				ID:          FirstString(m, "id", "product_id"),
				Name:        FirstString(m, "name", "title"),
				Description: FirstString(m, "description", "desc"),
				Cycle:       FirstString(m, "cycle", "billing_cycle"),
				Price:       FirstString(m, "price", "amount", "cost"),
			})
		}
	}
	return out, nil
}

// Contact is an entry of GET /api/contact.
type Contact struct {
	ID          string
	Email       string
	Type        string
	FirstName   string
	LastName    string
	CompanyName string
	PhoneNumber string
	Country     string
	State       string
	City        string
	Address     string
	Postcode    string
}

// ListContacts returns sub-account contacts (GET /api/contact).
func (c *Client) ListContacts(ctx context.Context) ([]Contact, error) {
	var raw any
	if err := c.Get(ctx, "/api/contact", nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "contacts")
	out := make([]Contact, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseContact(m))
		}
	}
	return out, nil
}

// GetContact returns contact details (GET /api/contact/:id).
func (c *Client) GetContact(ctx context.Context, id string) (*Contact, error) {
	var raw any
	if err := c.Get(ctx, "/api/contact/"+PathEscape(id), nil, &raw); err != nil {
		return nil, err
	}
	m, ok := AsMap(raw)
	if !ok {
		if list := ExtractList(raw, "contacts"); len(list) > 0 {
			m, _ = AsMap(list[0])
		}
	}
	if m == nil {
		return nil, fmt.Errorf("unexpected response shape for contact %s", id)
	}
	ct := parseContact(m)
	if ct.ID == "" {
		ct.ID = id
	}
	return &ct, nil
}

func parseContact(m map[string]any) Contact {
	return Contact{
		ID:          FirstString(m, "id", "contact_id"),
		Email:       FirstString(m, "email", "email_address"),
		Type:        FirstString(m, "type", "account_type"),
		FirstName:   FirstString(m, "firstname", "first_name"),
		LastName:    FirstString(m, "lastname", "last_name"),
		CompanyName: FirstString(m, "companyname", "company_name", "company"),
		PhoneNumber: FirstString(m, "phonenumber", "phone_number", "phone"),
		Country:     FirstString(m, "country"),
		State:       FirstString(m, "state"),
		City:        FirstString(m, "city"),
		Address:     FirstString(m, "address1", "address"),
		Postcode:    FirstString(m, "postcode", "zip"),
	}
}

// Certificate is an entry of GET /api/certificate.
type Certificate struct {
	ID      string
	Domain  string
	Status  string
	Type    string
	Expires string
}

// ListCertificates returns SSL certificates under the account
// (GET /api/certificate).
func (c *Client) ListCertificates(ctx context.Context) ([]Certificate, error) {
	var raw any
	if err := c.Get(ctx, "/api/certificate", nil, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "certificates", "services")
	out := make([]Certificate, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseCertificate(m))
		}
	}
	return out, nil
}

// GetCertificate returns certificate details (GET /api/certificate/:id).
func (c *Client) GetCertificate(ctx context.Context, id string) (*Certificate, error) {
	var raw any
	if err := c.Get(ctx, "/api/certificate/"+PathEscape(id), nil, &raw); err != nil {
		return nil, err
	}
	m, ok := AsMap(raw)
	if !ok {
		return nil, fmt.Errorf("unexpected response shape for certificate %s", id)
	}
	cert := parseCertificate(m)
	if cert.ID == "" {
		cert.ID = id
	}
	return &cert, nil
}

func parseCertificate(m map[string]any) Certificate {
	return Certificate{
		ID:      FirstString(m, "id", "certificate_id", "service_id"),
		Domain:  FirstString(m, "domain", "domain_name", "csr", "name"),
		Status:  FirstString(m, "status", "state"),
		Type:    FirstString(m, "type", "product", "product_name"),
		Expires: FirstString(m, "expires", "expiry", "expiry_date", "expires_date"),
	}
}

// URLLink is a short link managed by POST /api/url-shortener/shorten.
type URLLink struct {
	ID     string
	URL    string
	Short  string
	Label  string
	Clicks int64
}

// ShortenURL creates a short link (POST /api/url-shortener/shorten).
func (c *Client) ShortenURL(ctx context.Context, longURL, label string) (*URLLink, error) {
	var raw any
	if err := c.Post(ctx, "/api/url-shortener/shorten", Query("url", longURL, "label", label), &raw); err != nil {
		return nil, err
	}
	m, _ := AsMap(raw)
	if m == nil {
		m = map[string]any{}
	}
	link := parseURLLink(m)
	if link.URL == "" {
		link.URL = longURL
	}
	if link.Label == "" {
		link.Label = label
	}
	if link.ID == "" {
		return nil, fmt.Errorf("shorten response did not include a link id")
	}
	return &link, nil
}

// GetURLLink fetches a short link (GET /api/url-shortener/links/:id).
func (c *Client) GetURLLink(ctx context.Context, id string) (*URLLink, error) {
	var raw any
	if err := c.Get(ctx, "/api/url-shortener/links/"+PathEscape(id), nil, &raw); err != nil {
		return nil, err
	}
	m, ok := AsMap(raw)
	if !ok {
		if list := ExtractList(raw, "links"); len(list) > 0 {
			m, _ = AsMap(list[0])
		}
	}
	if m == nil {
		return nil, fmt.Errorf("unexpected response shape for link %s", id)
	}
	link := parseURLLink(m)
	if link.ID == "" {
		link.ID = id
	}
	return &link, nil
}

// DeleteURLLink removes a short link (DELETE /api/url-shortener/links/:id).
func (c *Client) DeleteURLLink(ctx context.Context, id string) error {
	return c.Delete(ctx, "/api/url-shortener/links/"+PathEscape(id), nil, nil)
}

// ListURLLinks lists the authenticated customer's short links
// (GET /api/url-shortener/links).
func (c *Client) ListURLLinks(ctx context.Context, page, perPage int64, search string) ([]URLLink, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", fmt.Sprintf("%d", page))
	}
	if perPage > 0 {
		q.Set("per_page", fmt.Sprintf("%d", perPage))
	}
	if search != "" {
		q.Set("search", search)
	}
	var raw any
	if err := c.Get(ctx, "/api/url-shortener/links", q, &raw); err != nil {
		return nil, err
	}
	list := ExtractList(raw, "links")
	out := make([]URLLink, 0, len(list))
	for _, item := range list {
		if m, ok := AsMap(item); ok {
			out = append(out, parseURLLink(m))
		}
	}
	return out, nil
}

func parseURLLink(m map[string]any) URLLink {
	return URLLink{
		ID:     FirstString(m, "id", "link_id"),
		URL:    FirstString(m, "url", "long_url", "original", "original_url"),
		Short:  FirstString(m, "short", "short_url", "shorturl", "link", "short_link"),
		Label:  FirstString(m, "label", "note"),
		Clicks: FirstInt64(m, "clicks", "hits", "visits", "click_count"),
	}
}
