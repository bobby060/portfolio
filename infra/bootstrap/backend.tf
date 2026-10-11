# This connects shared state with its stg backend
# When bootstrapping, comment this out for first local apply
# Once storage accounts exist and you have permission to write,
# Uncomment this and run `terraform init`
terraform {
  backend "azurerm" {
    resource_group_name = "rgp-portfolio-shared"
    container_name      = "tf-state-shared"
    key                 = "shared.terraform.tfstate"

    storage_account_name = "stgportfolio1542state"

    # Recommended security/auth settings
    use_azuread_auth = true # Use Microsoft Entra ID instead of access keys
    use_cli          = true # Allows local auth via 'az login'
  }
}