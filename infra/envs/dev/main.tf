

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
