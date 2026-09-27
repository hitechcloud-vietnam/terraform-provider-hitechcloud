# Examples

Runnable Terraform configurations for the HiTechCloud provider.

| Directory | Description |
| --- | --- |
| [`provider/`](provider) | Minimal provider configuration |
| [`resources/hitechcloud_dns_zone/`](resources/hitechcloud_dns_zone) | DNS zone + records |
| [`resources/hitechcloud_vm/`](resources/hitechcloud_vm) | VM, interface and firewall rule |
| [`resources/hitechcloud_ai_instance/`](resources/hitechcloud_ai_instance) | AI Factory SSH key, volume and GPU instance |
| [`resources/hitechcloud_s3_bucket/`](resources/hitechcloud_s3_bucket) | S3 bucket and sub-user |
| [`data-sources/hitechcloud_account/`](data-sources/hitechcloud_account) | Account, services and catalog lookups |

## Running

```shell
export HITECHCLOUD_TOKEN="your-api-token"
cd resources/hitechcloud_dns_zone
terraform init
terraform plan
```

All examples read the API token from the `HITECHCLOUD_TOKEN` environment
variable. Never hard-code credentials in these files.
