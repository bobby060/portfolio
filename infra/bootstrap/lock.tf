# Deleting this account would orphan every stack's state
resource "azurerm_management_lock" "tfstate" {
  name       = "tfstate-cannot-delete"
  scope      = azurerm_storage_account.tfstate.id
  lock_level = "CanNotDelete"
  notes      = "Holds all Terraform state"
}
