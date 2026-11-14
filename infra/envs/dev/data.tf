data "azurerm_resource_group" "dev" {
  name = "rgp-portfolio-${local.env}"
}

data "cloudflare_ip_ranges" "edge" {}

data "terraform_remote_state" "shared" {
  backend = "azurerm"
  config = {
    resource_group_name  = "rgp-portfolio-shared"
    storage_account_name = "stgportfolio1542state"
    container_name       = "tf-state-shared"
    key                  = "cloudflare.terraform.tfstate"

    use_azuread_auth = true
    use_cli          = true
  }
}