---
page_title: "hitechcloud_email_forwarding Resource - hitechcloud"
subcategory: "Domains"
description: |-
  Manages an email forwarding rule of a registered domain.
---

# hitechcloud_email_forwarding (Resource)

Manages an email forwarding rule of a registered domain
(`PUT /api/domain/{id}/emforwarding`). Destroying the resource clears the
forwarding rule.

## Example Usage

```terraform
resource "hitechcloud_email_forwarding" "info" {
  domain_id = "42"
  from      = "info@example.com"
  to        = "owner@example.net"
}
```

## Schema

### Required

- `domain_id` (String) ID of the registered domain. Changing this forces a new
  resource.
- `from` (String) Source email address to forward. Changing this forces a new
  resource.
- `to` (String) Destination email address receiving forwarded mail.

### Read-Only

- `id` (String) Composite identifier `domain_id/from`.

## Import

Import is supported using the composite ID `domain_id/from`:

```shell
terraform import hitechcloud_email_forwarding.info 42/info@example.com
```
