# Manage an AI Factory GPU instance together with an SSH key and a volume.

terraform {
  required_providers {
    hitechcloud = {
      source = "hitechcloud-vietnam/hitechcloud"
    }
  }
}

provider "hitechcloud" {}

resource "hitechcloud_ai_ssh_key" "main" {
  service_id = "1"
  name       = "laptop"
  public_key = file("~/.ssh/id_ed25519.pub")
  default    = true
}

resource "hitechcloud_ai_volume" "data" {
  service_id = "1"
  name       = "training-data"
  cloud      = "shade"
  region     = "hanoi-1"
  size_in_gb = 200
}

resource "hitechcloud_ai_instance" "gpu" {
  service_id          = "1"
  name                = "train-1"
  cloud               = "shade"
  region              = "hanoi-1"
  shade_instance_type = "A100"
  os                  = "ubuntu-22.04"

  ssh_key_id = hitechcloud_ai_ssh_key.main.key_id
  volume_ids = [hitechcloud_ai_volume.data.volume_id]

  tags = ["prod", "training"]
  envs = {
    MODE = "train"
  }

  auto_delete = false
  alert       = true
}

output "instance_id" {
  value = hitechcloud_ai_instance.gpu.instance_id
}
