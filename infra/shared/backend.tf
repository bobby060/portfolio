terraform {
  backend "azurerm" {
    resource_group_name  = "rgp-portfolio-shared"
    storage_account_name = "stgportfolio1542state"
    container_name       = "tf-state-shared"
    key                  = "cloudflare.terraform.tfstate"

    use_azuread_auth = true
    use_cli          = true
  }
}
