---
page_title: "hitechcloud_ai_template Resource - hitechcloud"
subcategory: "AI Factory"
description: |-
  Manages a template of the HiTechCloud AI Factory.
---

# hitechcloud_ai_template (Resource)

Manages a template of the HiTechCloud AI Factory
(`/api/service/{service_id}/templates`). `name`, `description` and
`is_public` are updated in place; content changes force a new resource.

## Example Usage

```terraform
resource "hitechcloud_ai_template" "pytorch" {
  service_id  = "1"
  name        = "pytorch-training"
  description = "PyTorch training environment"
  is_public   = false
}
```

## Schema

### Required

- `name` (String) Template name.
- `service_id` (String) ID of the hosted service. Changing this forces a new
  resource.

### Optional

- `description` (String) Template description.
- `is_public` (Boolean) Whether the template is public.
- `template` (String) Template content / definition. Changing this forces a new
  resource.

### Read-Only

- `id` (String) Composite identifier `service_id/template_id`.
- `template_id` (String) Identifier of the template as returned by the API.

## Import

Import is supported using the composite ID `service_id/template_id`:

```shell
terraform import hitechcloud_ai_template.pytorch 1/tpl-123
```
