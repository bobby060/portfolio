
resource "azurerm_storage_account" "tfstate" {

  name                     = "stgportfolio${local.stgid}state"
  resource_group_name      = azurerm_resource_group.shared.name
  location                 = azurerm_resource_group.shared.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  tags = local.tags

  shared_access_key_enabled       = false
  allow_nested_items_to_be_public = false
  min_tls_version                 = "TLS1_2"

  blob_properties {
    versioning_enabled = true

    delete_retention_policy { days = 30 }
    container_delete_retention_policy { days = 30 }
  }

}

resource "azurerm_storage_container" "tfstate" {
  for_each = toset(local.envs)

  name                  = "tf-state-${each.key}"
  storage_account_id    = azurerm_storage_account.tfstate.id
  container_access_type = "private"
}

resource "azurerm_storage_container" "shared-tfstate" {
  name                  = "tf-state-shared"
  storage_account_id    = azurerm_storage_account.tfstate.id
  container_access_type = "private"
}
