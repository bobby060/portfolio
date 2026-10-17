

module "container_app" {
  source = "../../container-app"

  resource_group_name = data.azurerm_resource_group.dev.name
  tags                = local.tags
  name                = local.name

  containername = local.containername
  containerpath = local.containerpath

  ingress_allowed_cidrs = data.cloudflare_ip_ranges.edge.ipv4_cidrs
}


resource "azurerm_container_app_custom_domain" "dev" {
  name                                     = var.custom_domain
  container_app_id                         = module.container_app.container_app_id
  container_app_environment_certificate_id = azurerm_container_app_environment_certificate.origin.id
  certificate_binding_type                 = "SniEnabled"
}



# ---- Access: who may reach the app ----
# Only effective while the hostname is proxied (dns_proxied = true).

resource "cloudflare_zero_trust_access_policy" "allowed_users" {
  account_id = var.cloudflare_account_id
  name       = "${local.name}-allowed-users"
  decision   = "allow"

  include = [for email in var.access_allowed_emails : { email = { email = email } }]
}

resource "cloudflare_zero_trust_access_application" "dev" {
  account_id       = var.cloudflare_account_id
  name             = local.name
  type             = "self_hosted"
  domain           = var.custom_domain
  session_duration = "24h"

  policies = [{
    id         = cloudflare_zero_trust_access_policy.allowed_users.id
    precedence = 1
  }]
}
