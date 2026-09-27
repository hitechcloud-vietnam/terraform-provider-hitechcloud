# Manage a virtual machine with a network interface and a firewall rule.

terraform {
  required_providers {
    hitechcloud = {
      source = "hitechcloud-vietnam/hitechcloud"
    }
  }
}

provider "hitechcloud" {}

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

variable "vm_password" {
  type      = string
  sensitive = true
}

resource "hitechcloud_vm_interface" "eth0" {
  service_id = "1"
  vm_id      = hitechcloud_vm.web.vm_id
  firewall   = true
}

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

output "vm_id" {
  value = hitechcloud_vm.web.vm_id
}
