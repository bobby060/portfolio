# This is the secret that authenticates with container
resource "random_password" "origin_auth" {
  for_each = local.hosts
  length   = 48
  special  = false # goes in an HTTP header; keep it to [A-Za-z0-9]
}

# Ruleset is applied to entire dns zone, which contains both target environments
resource "cloudflare_ruleset" "origin_auth" {
  zone_id = var.cloudflare_zone_id
  name    = "origin-auth"
  kind    = "zone"
  phase   = "http_request_late_transform"

  rules = [for env, host in local.hosts : {
    description = "origin auth: ${env}"
    expression  = "(http.host eq \"${host}\")"
    action      = "rewrite"
    enabled     = true
    action_parameters = {
      headers = {
        "X-Origin-Auth" = {
          operation = "set"
          value     = random_password.origin_auth[env].result
        }
      }
    }
  }]
}
