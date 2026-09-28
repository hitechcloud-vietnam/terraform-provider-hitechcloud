---
page_title: "hitechcloud_product_config Data Source - hitechcloud"
subcategory: "Billing"
description: |-
  Reads the order form configuration of one product.
---

# hitechcloud_product_config (Data Source)

Reads the order form configuration of one product (`GET /api/order/{product_id}`).

## Example Usage

```terraform
data "hitechcloud_product_config" "vm" {
  product_id = "1"
}

output "order_form" {
  value = data.hitechcloud_product_config.vm.config
}
```

## Schema

### Required

- `product_id` (String) Product identifier from `hitechcloud_products`.

### Read-Only

- `config` (Map of String) Order form fields and available options of the
  product.
- `id` (String) The product ID (same as `product_id`).
