locals {
  repo_name = split("/", local.repo)[1]
}

data "github_user" "me" {
  username = "bobby060"
}

# Environments gate which jobs can get each deploy identity's token.
resource "github_repository_environment" "env" {
  for_each            = toset(local.envs)
  repository          = local.repo_name
  environment         = each.key
  prevent_self_review = false

  dynamic "reviewers" {
    for_each = each.key == "prod" ? [1] : []
    content {
      users = [data.github_user.me.id]
    }
  }

  deployment_branch_policy {
    protected_branches     = false
    custom_branch_policies = true
  }
}

resource "github_repository_environment_deployment_policy" "main_only" {
  repository     = local.repo_name
  environment    = "prod"
  branch_pattern = "main"
}

resource "github_repository_environment_deployment_policy" "dev_only" {
  repository     = local.repo_name
  environment    = "dev"
  branch_pattern = "dev"
}

resource "github_actions_environment_variable" "azure_client_id" {
  for_each      = github_repository_environment.env
  repository    = local.repo_name
  environment   = each.value.environment
  variable_name = "AZURE_CLIENT_ID"
  value         = azurerm_user_assigned_identity.gh["deploy-${each.key}"].client_id
}

resource "github_actions_variable" "repo" {
  for_each = {
    AZURE_TENANT_ID       = data.azurerm_client_config.current.tenant_id
    AZURE_SUBSCRIPTION_ID = data.azurerm_client_config.current.subscription_id
    AZURE_PLAN_CLIENT_ID  = azurerm_user_assigned_identity.gh["plan"].client_id
  }
  repository    = local.repo_name
  variable_name = each.key
  value         = each.value
}
