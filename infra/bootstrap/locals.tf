locals {
  tags = {
    Project   = "Portfolio"
    ManagedBy = "Terraform"
  }

  envs = ["dev", "prod"]

  region = "eastus"

  # USed to make a unique stg account
  stgid = "1542"

}