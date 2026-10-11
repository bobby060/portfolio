resource "azurerm_log_analytics_workspace" "this" {
  name                = "${var.name}-law"
  location            = local.region
  resource_group_name = var.resource_group_name
  sku                 = "PerGB2018"
  retention_in_days   = 30

  tags = var.tags
}

resource "azurerm_container_app_environment" "this" {
  name                       = "${var.name}-app-env"
  location                   = local.region
  resource_group_name        = var.resource_group_name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.this.id

  tags = var.tags
}

resource "azurerm_container_app" "this" {
  name                         = "${var.name}-app"
  container_app_environment_id = azurerm_container_app_environment.this.id
  resource_group_name          = var.resource_group_name
  revision_mode                = "Single"

  template {
    container {
      name   = var.containername
      image  = var.container
      cpu    = 0.25
      memory = "0.5Gi"

    }
    min_replicas = 0
    max_replicas = 1

    http_scale_rule {
        name = "${var.name}-http-scale-rule"
        concurrent_requests = "100" #triggers additional scale up
    }
  }
}
