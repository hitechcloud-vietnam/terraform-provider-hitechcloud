---
page_title: "hitechcloud_partner_lead Resource - hitechcloud"
subcategory: "Portal"
description: |-
  Registers a partner lead / sales opportunity.
---

# hitechcloud_partner_lead (Resource)

Registers a partner lead / sales opportunity (`POST /api/partner/leads`). Leads
are immutable registrations: changing any attribute forces a new resource.

Duplicate emails of active leads are rejected by the API. Deleting the resource
removes it from Terraform state only and emits a warning.

## Example Usage

```terraform
resource "hitechcloud_partner_lead" "acme" {
  email        = "it@acme.example"
  company      = "ACME Corp"
  contact_name = "Nguyen Van A"
  phone        = "+84900000000"
}
```

## Schema

### Required

- `company` (String) Company name of the lead. Changing this forces a new
  resource.
- `email` (String) Lead's email address; duplicates with an active lead are
  rejected by the API. Changing this forces a new resource.

### Optional

- `contact_name` (String) Contact person's name. Changing this forces a new
  resource.
- `phone` (String) Contact phone number. Changing this forces a new resource.

### Read-Only

- `created_at` (String) Creation timestamp reported by the API.
- `id` (String) Lead identifier returned by the API.
- `status` (String) Current lead status reported by the API.

## Import

Import is supported using the lead ID (or its email):

```shell
terraform import hitechcloud_partner_lead.acme lead-1001
```
