---
page_title: "hitechcloud_ai_cluster Resource - hitechcloud"
subcategory: "AI Factory"
description: |-
  Manages a cluster of the HiTechCloud AI Factory.
---

# hitechcloud_ai_cluster (Resource)

Manages a cluster of the HiTechCloud AI Factory
(`/api/service/{service_id}/ai/cluster`). All attributes force a new resource.

## Example Usage

```terraform
resource "hitechcloud_ai_cluster" "shared" {
  service_id   = "1"
  name         = "shared-gpu"
  cloud        = "shade"
  region       = "hanoi-1"
  cluster_type = "A100x4"
  num_instances = 2

  ssh_key_id = hitechcloud_ai_ssh_key.main.key_id
}
```

## Schema

### Required

- `cluster_type` (String) Cluster type. Changing this forces a new resource.
- `cloud` (String) Cloud of the cluster. Changing this forces a new resource.
- `name` (String) Cluster name. Changing this forces a new resource.
- `num_instances` (Number) Number of instances in the cluster. Changing this
  forces a new resource.
- `region` (String) Region of the cluster. Changing this forces a new resource.
- `service_id` (String) ID of the hosted service. Changing this forces a new
  resource.

### Optional

- `os` (String) Operating system image. Changing this forces a new resource.
- `ssh_key_id` (String) SSH key attached to the cluster. Defaults to the
  account default key. Changing this forces a new resource.

### Read-Only

- `cluster_id` (String) Identifier of the cluster as returned by the API.
- `id` (String) Composite identifier `service_id/cluster_id`.
- `status` (String) Current cluster status.

## Import

Import is supported using the composite ID `service_id/cluster_id`:

```shell
terraform import hitechcloud_ai_cluster.shared 1/cl-123
```
