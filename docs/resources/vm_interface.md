---
page_title: "hitechcloud_vm_interface Resource - hitechcloud"
subcategory: "Compute"
description: |-
  Manages a network interface of a HiTechCloud virtual machine.
---

# hitechcloud_vm_interface (Resource)

Manages a network interface of a HiTechCloud virtual machine
(`/api/service/{service_id}/vm/{vm_id}/interface`).

## Example Usage

```terraform
resource "hitechcloud_vm_interface" "eth0" {
  service_id = "1"
  vm_id      = hitechcloud_vm.web.vm_id
  firewall   = true

  ipv4 = ["203.0.113.10"]
}
```

## Schema

### Required

- `service_id` (String) ID of the hosted service. Changing this forces a new
  resource.
- `vm_id` (String) ID of the VM. Changing this forces a new resource.

### Optional

- `bridge` (String) Network bridge. Defaults to the API default. Changing this
  forces a new resource.
- `firewall` (Boolean) Whether the firewall is enabled on the interface.
- `ipv4` (Set of String) IPv4 addresses attached to the interface.
- `ipv6` (Set of String) IPv6 addresses attached to the interface.

### Read-Only

- `id` (String) Composite identifier `service_id/vm_id/interface_id`.
- `interface_id` (String) Identifier of the interface as returned by the API.

## Import

Import is supported using the composite ID `service_id/vm_id/interface_id`:

```shell
terraform import hitechcloud_vm_interface.eth0 1/vm-123/if-456
```
