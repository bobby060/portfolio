# Full (strict): Cloudflare checks the origin cert against our Origin CA chain.
resource "cloudflare_zone_setting" "ssl" {
  zone_id    = var.cloudflare_zone_id
  setting_id = "ssl"
  value      = "strict"
}

import {
  to = cloudflare_zone_setting.ssl
  id = "${var.cloudflare_zone_id}/ssl"
}
