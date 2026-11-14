locals {
  tags = {
    Project   = "Portfolio"
    ManagedBy = "Terraform"
  }

  envs = ["dev", "prod"]

  region = "eastus"

  # USed to make a unique stg account
  stgid = "1542"

  repo = "bobby060/portfolio"

  gh_sub_prefix = "repo:bobby060@21338234/portfolio@1414640546"

  gh_identities = {
    plan        = "${local.gh_sub_prefix}:pull_request"
    deploy-dev  = "${local.gh_sub_prefix}:environment:dev"
    deploy-prod = "${local.gh_sub_prefix}:environment:prod"
  }


}