# Identities GitHub Actions log in as
resource "azurerm_user_assigned_identity" "gh" {
  for_each            = local.gh_identities
  name                = "id-gh-${each.key}"
  resource_group_name = azurerm_resource_group.shared.name
  location            = local.region
  tags                = local.tags
}

resource "azurerm_federated_identity_credential" "gh" {
  for_each                  = local.gh_identities
  name                      = "github-${each.key}"
  user_assigned_identity_id = azurerm_user_assigned_identity.gh[each.key].id
  audience                  = ["api://AzureADTokenExchange"]
  issuer                    = "https://token.actions.githubusercontent.com"
  subject                   = each.value
}


############################
## deploy-<env> scope
############################
resource "azurerm_role_assignment" "deploy_contributor" {
  for_each             = toset(local.envs)
  scope                = azurerm_resource_group.rgp[each.key].id
  principal_id         = azurerm_user_assigned_identity.gh["deploy-${each.key}"].principal_id
  role_definition_name = "Contributor"
  description          = "Update resources in own resource group"
}

resource "azurerm_role_assignment" "deploy_stg_contrib" {
  for_each             = toset(local.envs)
  scope                = azurerm_storage_container.tfstate[each.key].id
  principal_id         = azurerm_user_assigned_identity.gh["deploy-${each.key}"].principal_id
  role_definition_name = "Storage Blob Data Contributor"
  description          = "Read and write own state"
}

resource "azurerm_role_assignment" "deploy_reader" {
  for_each             = toset(local.envs)
  scope                = azurerm_storage_container.shared-tfstate.id
  principal_id         = azurerm_user_assigned_identity.gh["deploy-${each.key}"].principal_id
  role_definition_name = "Storage Blob Data Reader"
  description          = "Read shared state for origin secret"
}

############################
## plan
############################
# Can read env resources and env terraform states - cannot see shared
resource "azurerm_role_assignment" "plan_reader" {
  for_each             = toset(local.envs)
  scope                = azurerm_resource_group.rgp[each.key].id
  principal_id         = azurerm_user_assigned_identity.gh["plan"].principal_id
  role_definition_name = "Reader"
  description          = "View all resources"
}

resource "azurerm_role_assignment" "plan_stg_reader" {
  for_each             = toset(local.envs)
  scope                = azurerm_storage_container.tfstate[each.key].id
  principal_id         = azurerm_user_assigned_identity.gh["plan"].principal_id
  role_definition_name = "Storage Blob Data Reader"
  description          = "Read env states"
}

resource "azurerm_role_assignment" "plan_stg_reader_shared" {
  scope                = azurerm_storage_container.shared-tfstate.id
  principal_id         = azurerm_user_assigned_identity.gh["plan"].principal_id
  role_definition_name = "Storage Blob Data Reader"
  description          = "Read shared states"
}


