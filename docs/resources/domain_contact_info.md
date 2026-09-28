---
page_title: "hitechcloud_domain_contact_info Resource - hitechcloud"
subcategory: "Domains"
description: |-
  Manages the contact information of a registered domain.
---

# hitechcloud_domain_contact_info (Resource)

Manages the contact information of a registered domain
(`PUT /api/domain/{id}/contact`). The `contact_info` value is passed verbatim to
the API as the `contact_info` query parameter. Destroying the resource leaves
the contacts unchanged (the API has no reset).

## Example Usage

```terraform
resource "hitechcloud_domain_contact_info" "example" {
  domain_id    = "42"
  contact_info = "registrant=10&admin=11&tech=12&billing=13"
}
```

## Schema

### Required

- `contact_info` (String) Contact information payload accepted by
  `PUT /api/domain/{id}/contact` (the `contact_info` parameter).
- `domain_id` (String) ID of the registered domain. Changing this forces a new
  resource.

### Read-Only

- `id` (String) The domain ID (same as `domain_id`).

## Import

Import is supported using the domain ID:

```shell
terraform import hitechcloud_domain_contact_info.example 42
```
