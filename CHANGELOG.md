# Changelog

## 0.1.0 (unreleased)

FEATURES:

* **New resource:** `hitechcloud_dns_zone`
* **New resource:** `hitechcloud_dns_record`
* **New resource:** `hitechcloud_domain_dns_record`
* **New resource:** `hitechcloud_domain_settings`
* **New resource:** `hitechcloud_vm`
* **New resource:** `hitechcloud_vm_interface`
* **New resource:** `hitechcloud_vm_firewall_rule`
* **New resource:** `hitechcloud_service_ip`
* **New resource:** `hitechcloud_ai_instance`
* **New resource:** `hitechcloud_ai_ssh_key`
* **New resource:** `hitechcloud_ai_volume`
* **New resource:** `hitechcloud_ai_template`
* **New resource:** `hitechcloud_ai_cluster`
* **New resource:** `hitechcloud_s3_bucket`
* **New resource:** `hitechcloud_s3_subuser`
* **New resource:** `hitechcloud_url_shortener_link`
* **New data source:** `hitechcloud_account`
* **New data source:** `hitechcloud_contacts`
* **New data source:** `hitechcloud_services`
* **New data source:** `hitechcloud_service`
* **New data source:** `hitechcloud_invoices`
* **New data source:** `hitechcloud_payment_methods`
* **New data source:** `hitechcloud_products`
* **New data source:** `hitechcloud_certificates`
* **New data source:** `hitechcloud_domains`
* **New data source:** `hitechcloud_domain`
* **New data source:** `hitechcloud_domain_tlds`
* **New data source:** `hitechcloud_dns_zones`
* **New data source:** `hitechcloud_vms`
* **New data source:** `hitechcloud_ai_instance_types`
* **New data source:** `hitechcloud_ai_ssh_keys`
* **New data source:** `hitechcloud_ai_volumes`
* **New data source:** `hitechcloud_ai_cluster_types`
* **New data source:** `hitechcloud_s3_buckets`
* **New data source:** `hitechcloud_s3_subusers`

NOTES:

* Initial release targeting the HiTechCloud User API
  (`https://api.hitechcloud.vn`) with bearer-token authentication via the
  `token` attribute or the `HITECHCLOUD_TOKEN` environment variable.
* Built with the Terraform Plugin Framework (protocol 6).
