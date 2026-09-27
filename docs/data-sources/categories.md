---
page_title: "hitechcloud_categories Data Source - hitechcloud"
subcategory: "Billing"
description: |-
  Lists the product categories of the HiTechCloud catalog.
---

# hitechcloud_categories (Data Source)

Lists the product categories of the HiTechCloud catalog (`GET /api/category`).
Useful to discover the `category_id` used by the
[`hitechcloud_products`](products.md) data source.

## Example Usage

```terraform
data "hitechcloud_categories" "all" {}

output "category_names" {
  value = [for c in data.hitechcloud_categories.all.categories : c.name]
}
```

## Schema

### Read-Only

- `categories` (Attributes List) Product categories. (see
  [nested schema](#nestedatt--categories))
- `id` (String) Data source identifier (`categories`).

<a id="nestedatt--categories"></a>

### Nested Schema for `categories`

Read-Only:

- `description` (String) Category description.
- `id` (String) Category identifier.
- `name` (String) Category name.
