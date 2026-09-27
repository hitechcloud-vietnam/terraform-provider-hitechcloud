---
page_title: "hitechcloud_vm_firewall_rule Resource - hitechcloud"
subcategory: "Compute"
description: |-
  Manages a firewall rule of a HiTechCloud virtual machine.
---

# hitechcloud_vm_firewall_rule (Resource)

Manages a firewall rule of a HiTechCloud virtual machine
(`/api/service/{service_id}/vm/{vm_id}/firewall`). Rules are positional; the
`position` is assigned at creation and changing any attribute forces a new
resource.

## Example Usage

```terraform
resource "hitechcloud_vm_firewall_rule" "ssh" {
  service_id = "1"
  vm_id      = hitechcloud_vm.web.vm_id

  action        = "accept"
  type          = "in"
  protocol      = "tcp"
  address_start = "0.0.0.0"
  address_end   = "255.255.255.255"
  port_start    = 22
  port_end      = 22
  comment       = "SSH"
}
```

## Schema

### Required

- `action` (String) Rule action, for example `accept` or `drop`. Changing this
  forces a new resource.
- `address_end` (String) End of the remote address range. Changing this forces
  a new resource.
- `address_start` (String) Start of the remote address range. Changing this
  forces a new resource.
- `service_id` (String) ID of the hosted service. Changing this forces a new
  resource.
- `type` (String) Traffic direction (`in`/`out`). Changing this forces a new
  resource.
- `vm_id` (String) ID of the VM. Changing this forces a new resource.

### Optional

- `comment` (String) Rule comment. Changing this forces a new resource.
- `port_end` (Number) End of the port range. Changing this forces a new
  resource.
- `port_start` (Number) Start of the port range. Changing this forces a new
  resource.
- `protocol` (String) Protocol (`tcp`, `udp`, `icmp`, ...). Changing this
  forces a new resource.

### Read-Only

- `id` (String) Composite identifier `service_id/vm_id/position`.
- `position` (Number) Position of the rule in the firewall chain.

## Import

Import is supported using the composite ID `service_id/vm_id/position`:

```shell
terraform import hitechcloud_vm_firewall_rule.ssh 1/vm-123/3
```
