---
page_title: "hitechcloud_ai_instance Resource - hitechcloud"
subcategory: "AI Factory"
description: |-
  Manages a HiTechCloud AI Factory GPU instance.
---

# hitechcloud_ai_instance (Resource)

Manages a HiTechCloud AI Factory GPU instance
(`/api/service/{service_id}/ai/instance`). Creation waits until the instance
reaches a stable state; deletion waits until the instance is released.

## Example Usage

```terraform
resource "hitechcloud_ai_instance" "gpu" {
  service_id          = "1"
  name                = "train-1"
  cloud               = "shade"
  region              = "hanoi-1"
  shade_instance_type = "A100"
  os                  = "ubuntu-22.04"

  ssh_key_id = hitechcloud_ai_ssh_key.main.key_id

  tags = ["prod", "training"]
  envs = {
    MODE = "train"
  }

  auto_delete = false
  alert       = true
}
```

## Schema

### Required

- `cloud` (String) Cloud of the instance. Changing this forces a new resource.
- `name` (String) Instance name.
- `region` (String) Region of the instance. Changing this forces a new
  resource.
- `service_id` (String) ID of the hosted service. Changing this forces a new
  resource.
- `shade_instance_type` (String) Instance type. Changing this forces a new
  resource.

### Optional

- `alert` (Boolean) Enable alerting for the instance.
- `auto_delete` (Boolean) Automatically delete the instance when it stops.
- `envs` (Map of String) Environment variables passed to the instance.
- `launch_configuration` (String) Launch configuration (JSON).
- `os` (String) Operating system image. Changing this forces a new resource.
- `shade_cloud` (Boolean) Whether this is a Shade cloud instance. Changing this
  forces a new resource.
- `ssh_key_id` (String) SSH key attached to the instance. Defaults to the
  account default key. Changing this forces a new resource.
- `tags` (Set of String) Instance tags.
- `template_id` (String) Template used to provision the instance. Defaults to
  the API default. Changing this forces a new resource.
- `volume_ids` (Set of String) Volumes attached at creation. Changing this
  forces a new resource.
- `volume_mount` (String) Volume mount configuration (JSON).

### Read-Only

- `id` (String) Composite identifier `service_id/instance_id`.
- `instance_id` (String) Identifier of the instance as returned by the API.
- `status` (String) Current instance status.

## Import

Import is supported using the composite ID `service_id/instance_id`:

```shell
terraform import hitechcloud_ai_instance.gpu 1/ai-123
```
