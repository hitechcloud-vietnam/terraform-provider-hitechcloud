---
page_title: "hitechcloud_url_shortener_links Data Source - hitechcloud"
subcategory: "Portal"
description: |-
  Lists the shortened URLs of the HiTechCloud portal.
---

# hitechcloud_url_shortener_links (Data Source)

Lists the shortened URLs of the HiTechCloud portal
(`GET /api/url-shortener/links`).

## Example Usage

```terraform
data "hitechcloud_url_shortener_links" "all" {
  per_page = 50
}

output "short_urls" {
  value = [for l in data.hitechcloud_url_shortener_links.all.links : l.short_url]
}
```

## Schema

### Optional

- `page` (Number) Page number. Defaults to `1`.
- `per_page` (Number) Results per page.
- `search` (String) Search term.

### Read-Only

- `id` (String) Data source identifier (`url_shortener_links`).
- `links` (Attributes List) Shortened links. (see
  [nested schema](#nestedatt--links))

<a id="nestedatt--links"></a>

### Nested Schema for `links`

Read-Only:

- `clicks` (Number) Click counter.
- `id` (String) Link identifier.
- `label` (String) Label / alias.
- `short_url` (String) Short URL.
- `url` (String) Target URL.
