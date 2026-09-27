---
page_title: "hitechcloud_ai_volume Resource - hitechcloud"
subcategory: "AI Factory"
description: |-
  Manages a volume of the HiTechCloud AI Factory.
---

# hitechcloud_ai_volume (Resource)

Manages a volume of the HiTechCloud AI Factory
(`/api/service/{service_id}/ai/volume`). All attributes force a new resource.

## Example Usage

```terraform
resource "hitechcloud_ai_volume" "data" {
  service_id = "1"
  name       = "training-data"
  cloud      = "shade"
  region     = "hanoi-1"
  size_in_gb = 200
}
```

## Schema

### Required

- `cloud` (String) Cloud of the volume. Changing this forces a new resource.
- `name` (String) Volume name. Changing this forces a new resource.
- `region` (String) Region of the volume. Changing this forces a new resource.
- `service_id` (String) ID of the hosted service. Changing this forces a new
  resource.
- `size_in_gb` (Number) Volume size in GB. Changing this forces a new resource.

### Read-Only

- `id` (String) Composite identifier `service_id/volume_id`.
- `status` (String) Current volume status.
- `volume_id` (String) Identifier of the volume as returned by the API.

## Import

Import is supported using the composite ID `service_id/volume_id`:

```shell
terraform import hitechcloud_ai_volume.data 1/vol-123
```
