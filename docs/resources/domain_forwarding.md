---
page_title: "hitechcloud_domain_forwarding Resource - hitechcloud"
subcategory: "Domains"
description: |-
  Manages URL forwarding (web redirect) of a registered domain.
---

# hitechcloud_domain_forwarding (Resource)

Manages URL forwarding (web redirect) of a registered domain
(`PUT /api/domain/{id}/forwarding`). Destroying the resource disables
forwarding.

## Example Usage

```terraform
resource "hitechcloud_domain_forwarding" "example" {
  domain_id = "42"
  url       = "https://www.example.com"
  frame     = false
}
```

## Schema

### Required

- `domain_id` (String) ID of the registered domain. Changing this forces a new
  resource.
- `url` (String) Destination URL the domain should redirect to.

### Optional

- `frame` (Bool) When true, the destination is rendered inside a frame instead of
  an HTTP redirect.

### Read-Only

- `id` (String) The domain ID (same as `domain_id`).

## Import

Import is supported using the domain ID:

```shell
terraform import hitechcloud_domain_forwarding.example 42
```
