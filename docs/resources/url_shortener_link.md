---
page_title: "hitechcloud_url_shortener_link Resource - hitechcloud"
subcategory: "Portal"
description: |-
  Manages a shortened URL in the HiTechCloud portal.
---

# hitechcloud_url_shortener_link (Resource)

Manages a shortened URL in the HiTechCloud portal
(`/api/url-shortener/shorten`). Changing `url` or `label` forces a new
resource.

## Example Usage

```terraform
resource "hitechcloud_url_shortener_link" "docs" {
  url   = "https://docs.example.com/getting-started"
  label = "getting-started"
}

output "short_url" {
  value = hitechcloud_url_shortener_link.docs.short_url
}
```

## Schema

### Required

- `url` (String) Target URL. Changing this forces a new resource.

### Optional

- `label` (String) Optional label / alias for the short link. Changing this
  forces a new resource.

### Read-Only

- `clicks` (Number) Number of clicks recorded on the link.
- `id` (String) Identifier of the link as returned by the API.
- `short_url` (String) The resulting short URL.

~> **Note:** This resource does not support import; existing short links must
be re-created.
