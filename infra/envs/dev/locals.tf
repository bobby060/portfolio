locals {
  env  = "dev"
  name = "portfolio-${local.env}"

  tags = {
    Project   = "Portfolio"
    ManagedBy = "Terraform"
    env       = local.env
  }

  containername = "site"
  containerpath = var.container_path
}
