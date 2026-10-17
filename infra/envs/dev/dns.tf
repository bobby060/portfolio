# ACA checks this TXT record to prove we own the custom domain.
resource "cloudflare_dns_record" "azure_verification" {
  zone_id = var.cloudflare_zone_id
  name    = "asuid.${var.custom_domain}"
  type    = "TXT"
  content = module.container_app.custom_domain_verification_id
  ttl     = 3600
}

resource "cloudflare_dns_record" "app" {
  zone_id = var.cloudflare_zone_id
  name    = var.custom_domain
  type    = "CNAME"
  content = module.container_app.fqdn
  ttl     = 1 # automatic when proxied
  proxied = var.dns_proxied

  depends_on = [cloudflare_dns_record.azure_verification]
}

# Full (strict): Cloudflare checks the origin cert against our Origin CA chain.
resource "cloudflare_zone_setting" "ssl" {
  zone_id    = var.cloudflare_zone_id
  setting_id = "ssl"
  value      = "strict"
}