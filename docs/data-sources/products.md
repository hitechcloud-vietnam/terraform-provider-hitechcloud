---
page_title: "hitechcloud_products Data Source - hitechcloud"
subcategory: "Billing"
description: |-
  Lists the purchasable products of a HiTechCloud product category.
---

# hitechcloud_products (Data Source)

Lists the purchasable products of a HiTechCloud product category
(`GET /api/category/{category_id}/product`).

## Example Usage

```terraform
data "hitechcloud_products" "hosting" {
  category_id = "5"
}

output "product_names" {
  value = [for p in data.hitechcloud_products.hosting.products : p.name]
}
```

## Schema

### Required

- `category_id` (String) Product category identifier.

### Read-Only

- `id` (String) Data source identifier (same as `category_id`).
- `products` (Attributes List) Products in the category. (see
  [nested schema](#nestedatt--products))

<a id="nestedatt--products"></a>

### Nested Schema for `products`

Read-Only:

- `cycle` (String) Billing cycle.
- `description` (String) Product description.
- `id` (String) Product identifier.
- `name` (String) Product name.
- `price` (String) Price.
