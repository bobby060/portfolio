
data "azurerm_client_config" "current" {}

resource "azurerm_role_assignment" "stg_data_contrib" {
  scope                = azurerm_storage_account.tfstate.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         =  data.azurerm_client_config.current.object_id
}