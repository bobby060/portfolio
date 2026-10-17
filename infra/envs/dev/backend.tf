# State lives in the bootstrap-created storage account, one container per env
# (infra/bootstrap/stg.tf). Auth is Entra ID via `az login`; no account keys.
terraform {
  backend "azurerm" {
    resource_group_name  = "rgp-portfolio-shared"
    storage_account_name = "stgportfolio1542state"
    container_name       = "tf-state-dev"
    key                  = "dev.terraform.tfstate"

    use_azuread_auth = true
    use_cli          = true
  }
}
