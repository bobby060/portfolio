data "azurerm_resource_group" "dev" {
  name = "rgp-portfolio-${local.env}"
}

data "cloudflare_ip_ranges" "edge" {}