# Changelog

## 1.2.0 (September 28, 2026)

FEATURES:

* **New resource:** `hitechcloud_service_label` — manage the display label of a
  service (`POST /api/service/{id}/label`).
* **New resource:** `hitechcloud_service_billing_cycle` — manage the billing
  cycle of a service (`POST /api/service/{id}/cycle`).
* **New resource:** `hitechcloud_pmg_domain` — protect a mail domain through a
  PMG service and set its mail transport (`POST /api/service/{id}/htcpmg/domains`,
  `/transport`).
* **New resource:** `hitechcloud_pve_rdns` — reverse DNS (PTR) of a Proxmox
  service IP (`POST /api/service/{id}/htcpve/rdns`).
* **New resource:** `hitechcloud_ipam_rdns` — reverse DNS (PTR) of an IPAM
  address (`POST /api/service/{id}/htcipam/rdns`).
* **New resource:** `hitechcloud_pve_backup` — take Proxmox VE backups
  (`POST /api/service/{id}/htcpve/backups`).
* **New resource:** `hitechcloud_vm_rdns` — reverse DNS (PTR) of a VM IP
  (`POST /api/service/{id}/vms/{vmid}/rdns`).
* **New resource:** `hitechcloud_partner_payout` — request partner wallet
  payouts (`POST /api/partner/payouts`).
* **New resource:** `hitechcloud_affiliate_adv_voucher` — manage advanced
  affiliate discount vouchers (`POST/DELETE /api/affiliates_adv/{client_id}/vouchers`).
* **New resource:** `hitechcloud_domain_contact_info` — manage the contact
  information of a registered domain (`PUT /api/domain/{id}/contact`).
* **New data source:** `hitechcloud_affiliate` — affiliate summary, campaigns,
  commissions, payouts, vouchers and commission plans.
* **New data source:** `hitechcloud_affiliate_adv` — advanced affiliate profile,
  statistics, referrals, vouchers and commissions.
* **New data source:** `hitechcloud_locations` — supported countries and states.
* **New data source:** `hitechcloud_server_stock` — dedicated server stock.
* **New data source:** `hitechcloud_service_resources` — service resource
  summary and upgrade options.
* **New data source:** `hitechcloud_news` — news and knowledgebase categories.
* **New data source:** `hitechcloud_product_config` — order form configuration
  of a product.
* **New data source:** `hitechcloud_mfa_status` — passkey and email MFA status
  of a user.

NOTES:

* The client now covers **348/348** Postman collection endpoints, including the
  newly added `GET /api/dns`, `GET /api/service/{id}/ip`,
  `GET /api/service/{id}/cluster` and filtered `GET /api/statuses` wrappers.
* Terraform surface: 35 resources and 47 data sources.

SECURITY:

* Upgraded `golang.org/x/crypto` to v0.52.0 (13 Dependabot alerts, 7 critical).
* Upgraded `golang.org/x/net` to v0.55.0 (1 Dependabot alert).
* Upgraded `google.golang.org/grpc` to v1.83.1 (4 Dependabot alerts).
* Raised the module's minimum Go version to 1.26 and aligned CI to Go 1.26.x.

## 1.1.0 (September 28, 2026)

FEATURES:

* **New resource:** `hitechcloud_dnssec_key` — manage DNSSEC (DS) keys of a
  registered domain (`PUT/DELETE /api/domain/{id}/dnssec`).
* **New resource:** `hitechcloud_email_forwarding` — manage domain email
  forwarding rules (`PUT /api/domain/{id}/emforwarding`).
* **New resource:** `hitechcloud_domain_forwarding` — manage URL forwarding of
  a registered domain (`PUT /api/domain/{id}/forwarding`).
* **New resource:** `hitechcloud_service_autorenew` — manage automatic renewal
  of services and domains (`PUT /api/willexpired/{type}/{id}/autorenew`).
* **New resource:** `hitechcloud_vm_snapshot` — take Proxmox VM snapshots
  (`POST /api/service/{id}/htcpve/snapshots`).
* **New resource:** `hitechcloud_partner_lead` — register partner leads
  (`POST /api/partner/leads`).
* **New data source:** `hitechcloud_dnssec_keys` — list DNSSEC keys and
  available flags of a domain.
* **New data source:** `hitechcloud_domain_contact` — domain contacts, EPP
  code, registrar lock and ID protection state.
* **New data source:** `hitechcloud_pbs` — Proxmox Backup Server connection
  info, usage, metrics, snapshots and groups.
* **New data source:** `hitechcloud_s3_connection` — S3 connection info,
  credentials, usage and metrics.
* **New data source:** `hitechcloud_pve` — Proxmox VE status, VM inventory,
  IPs, backups, snapshots and usage.
* **New data source:** `hitechcloud_ipam` — IPAM IPs, subnets and reverse DNS.
* **New data source:** `hitechcloud_partner` — partner profile, tiers, rate
  card, wallet and recent leads.
* **New data source:** `hitechcloud_willexpired` — expiring services/domains,
  summary and upcoming invoices.
* **New data source:** `hitechcloud_account_logs` — account activity log.

## 1.0.0 (September 28, 2026)

FEATURES:

* **New provider authentication:** login via `username`/`password`
  (`POST /api/login`) returns two tokens (`token` + `refresh_token`); the
  access token is renewed automatically with `POST /api/token` and failed
  requests are replayed. `refresh_token` can also be provided directly.

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
* **New resource:** `hitechcloud_contact`
* **New resource:** `hitechcloud_ticket`
* **New resource:** `hitechcloud_rdns`
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
* **New data source:** `hitechcloud_balance`
* **New data source:** `hitechcloud_categories`
* **New data source:** `hitechcloud_payment_fees`
* **New data source:** `hitechcloud_whois`
* **New data source:** `hitechcloud_domain_dns_types`
* **New data source:** `hitechcloud_domain_availability`
* **New data source:** `hitechcloud_url_shortener_links`
* **New data source:** `hitechcloud_tickets`
* **New data source:** `hitechcloud_ticket_departments`
* **New data source:** `hitechcloud_notifications`
* **New data source:** `hitechcloud_statuses`

NOTES:

* Initial release targeting the HiTechCloud User API
  (`https://api.hitechcloud.vn`) with bearer-token authentication via the
  `token` attribute or the `HITECHCLOUD_TOKEN` environment variable, or
  username/password login (`HITECHCLOUD_USERNAME` / `HITECHCLOUD_PASSWORD`)
  with automatic token renewal.
* The Go client covers the full Postman surface: account, billing, support
  (tickets, news, knowledge base, notifications), DNS, domains (incl. DNSSEC,
  forwarding, EPP, renew/order), SSL certificates, compute/VM lifecycle,
  AI Factory, storage (S3, PBS), bare metal/colocation/Hosting, partner
  program, eKYC, MFA (passkey/email), WillExpired renewals and more.
* Built with the Terraform Plugin Framework (protocol 6).
* Releases are GPG-signed; the public signing key is kept in the
  environment/secret store only and is not committed to the repository.
