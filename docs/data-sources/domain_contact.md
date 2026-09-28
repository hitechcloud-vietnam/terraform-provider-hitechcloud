---
page_title: "hitechcloud_domain_contact Data Source - hitechcloud"
subcategory: "Domains"
description: |-
  Reads contact and transfer information of a registered domain.
---

# hitechcloud_domain_contact (Data Source)

Reads contact and transfer information of a registered domain
(`GET /api/domain/{id}/contact`, `GET /api/domain/{id}/epp`).

## Example Usage

```terraform
data "hitechcloud_domain_contact" "example" {
  domain_id = "42"
}

output "epp_code" {
  value     = data.hitechcloud_domain_contact.example.epp_code
  sensitive = true
}
```

## Schema

### Required

- `domain_id` (String) ID of the registered domain.

### Read-Only

- `admin_contact_id` (String) Admin contact ID.
- `billing_contact_id` (String) Billing contact ID.
- `epp_code` (String, Sensitive) EPP transfer code of the domain.
- `id` (String) The domain ID (same as `domain_id`).
- `id_protection` (Bool) Whether WHOIS ID protection is enabled.
- `registrar_lock` (Bool) Whether the domain is registrar-locked.
- `registrant_contact_id` (String) Registrant contact ID.
- `tech_contact_id` (String) Tech contact ID.
