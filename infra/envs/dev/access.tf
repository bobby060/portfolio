
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