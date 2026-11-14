terraform {
  required_version = ">= 1.5.0"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 5.9.0"
    }
    github = {
    source  = "integrations/github"
    version = "~> 6.0"
  }

  }
}

provider "azurerm" {
  features {}
  storage_use_azuread = true
}

provider "github" {
  owner = "bobby060"
}
