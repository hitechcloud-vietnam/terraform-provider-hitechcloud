---
page_title: "hitechcloud_service_autorenew Resource - hitechcloud"
subcategory: "Billing"
description: |-
  Manages the automatic-renewal flag of a service or domain.
---

# hitechcloud_service_autorenew (Resource)

Manages the automatic-renewal flag of a service or domain
(`PUT /api/willexpired/{type}/{id}/autorenew`). Destroying the resource turns
auto-renew off.

## Example Usage

```terraform
resource "hitechcloud_service_autorenew" "vm" {
  item_type  = "service"
  item_id    = "1"
  autorenew  = true
}

resource "hitechcloud_service_autorenew" "domain" {
  item_type = "domain"
  item_id   = "example.com"
  autorenew = true
}
```

## Schema

### Required

- `autorenew` (Bool) Whether automatic renewal is enabled for the item.
- `item_id` (String) Service ID (`hb_accounts.id`) or domain name
  (`hb_domains.id`). Changing this forces a new resource.
- `item_type` (String) Item type: `service` or `domain`. Changing this forces a
  new resource.

### Read-Only

- `id` (String) Composite identifier `item_type/item_id`.

## Import

Import is supported using the composite ID `item_type/item_id`:

```shell
terraform import hitechcloud_service_autorenew.vm service/1
```
