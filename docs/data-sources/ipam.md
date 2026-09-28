---
page_title: "hitechcloud_ipam Data Source - hitechcloud"
subcategory: "Network"
description: |-
  Reads an IPAM service: allocated IPs, subnets and reverse DNS records.
---

# hitechcloud_ipam (Data Source)

Reads an IPAM (IP address management) service: allocated IPs, subnets and reverse
DNS records (`GET /api/service/{id}/htcipam/ips`, `/htcipam/subnets`,
`/htcipam/rdns`).

## Example Usage

```terraform
data "hitechcloud_ipam" "pool" {
  service_id = "1"
}

output "allocated_ips" {
  value = [for ip in data.hitechcloud_ipam.pool.ips : ip.ip]
}
```

## Schema

### Required

- `service_id` (String) HiTechCloud service ID of the IPAM service.

### Read-Only

- `id` (String) The service ID (same as `service_id`).
- `ips` (Attributes List) IP addresses allocated from the IPAM pool. (see
  [nested schema](#nestedatt--ips))
- `rdns` (Attributes List) Reverse DNS records of the service. (see
  [nested schema](#nestedatt--rdns))
- `subnets` (Attributes List) Subnets managed by the IPAM service. (see
  [nested schema](#nestedatt--subnets))

<a id="nestedatt--ips"></a>

### Nested Schema for `ips`

Read-Only:

- `comment` (String) Operator comment.
- `ip` (String) IP address.
- `status` (String) Allocation status.
- `subnet` (String) Parent subnet.

<a id="nestedatt--rdns"></a>

### Nested Schema for `rdns`

Read-Only:

- `hostname` (String) PTR hostname.
- `ip` (String) IP address.

<a id="nestedatt--subnets"></a>

### Nested Schema for `subnets`

Read-Only:

- `capacity` (Number) Number of addresses in the subnet.
- `cidr` (String) Subnet in CIDR notation.
- `gateway` (String) Gateway address.
- `type` (String) Address family (ipv4/ipv6).
