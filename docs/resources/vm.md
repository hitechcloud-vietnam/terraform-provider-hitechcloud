---
page_title: "hitechcloud_vm Resource - hitechcloud"
subcategory: "Compute"
description: |-
  Manages a HiTechCloud virtual machine.
---

# hitechcloud_vm (Resource)

Manages a HiTechCloud virtual machine (`/api/service/{service_id}/vm`).
Creation waits until the VM leaves its transient provisioning states;
deletion waits until the VM is fully released.

## Example Usage

```terraform
resource "hitechcloud_vm" "web" {
  service_id  = "1"
  label       = "web-1"
  template_id = "ubuntu-22.04"
  hostname    = "web-1"
  memory      = 4096
  cpu         = 2
  disk        = 40

  password = var.vm_password
}
```

## Schema

### Required

- `label` (String) Display label of the VM.
- `service_id` (String) ID of the hosted service owning the VM. Changing this
  forces a new resource.

### Optional

- `cpu` (Number) Number of virtual CPUs.
- `cpu_share` (Number) CPU share (relative weight).
- `disk` (Number) Disk size in GB. Changing this forces a new resource.
- `hostname` (String) OS hostname.
- `license_key` (String) License key for licensed templates.
- `memory` (Number) Memory in MB.
- `note` (String) Free-form note.
- `password` (String, Sensitive) Root/Administrator password.
- `swap` (Number) Swap size in GB. Changing this forces a new resource.
- `template_id` (String) OS template identifier. Changing this forces a new
  resource.

### Read-Only

- `id` (String) Composite identifier `service_id/vm_id`.
- `ips` (Set of String) IP addresses of the VM.
- `status` (String) Current VM status.
- `vm_id` (String) Identifier of the VM as returned by the API.

## Import

Import is supported using the composite ID `service_id/vm_id`:

```shell
terraform import hitechcloud_vm.web 1/vm-123
```
