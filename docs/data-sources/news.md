---
page_title: "hitechcloud_news Data Source - hitechcloud"
subcategory: "Support"
description: |-
  Reads HiTechCloud news and knowledgebase categories.
---

# hitechcloud_news (Data Source)

Reads HiTechCloud news and knowledgebase categories
(`GET /api/news`, `GET /api/knowledgebase`).

## Example Usage

```terraform
data "hitechcloud_news" "latest" {}

output "headlines" {
  value = [for n in data.hitechcloud_news.latest.news : n.title]
}
```

## Schema

### Read-Only

- `id` (String) Static identifier for this query (`news`).
- `knowledgebase` (Attributes List) Knowledgebase categories. (see
  [nested schema](#nestedatt--knowledgebase))
- `news` (Attributes List) News articles. (see
  [nested schema](#nestedatt--news))

<a id="nestedatt--news"></a>

### Nested Schema for `news`

Read-Only:

- `id` (String) Article identifier.
- `title` (String) Article title.
- `date` (String) Publication date.
- `extra` (Map of String) Remaining fields of the row.

<a id="nestedatt--knowledgebase"></a>

### Nested Schema for `knowledgebase`

Read-Only:

- `id` (String) Category identifier.
- `title` (String) Category name.
- `extra` (Map of String) Remaining fields of the row.
