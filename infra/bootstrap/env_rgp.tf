resource "azurerm_resource_group" "shared" {

  name     = "rgp-portfolio-shared"
  location = local.region

  tags = local.tags
}

resource "azurerm_resource_group" "rgp" {

  for_each = toset(local.envs)

  name     = "rgp-portfolio-${each.key}"
  location = local.region

  tags = local.tags
}

