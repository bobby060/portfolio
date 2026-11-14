locals {
  tags = {
    Project   = "Portfolio"
    ManagedBy = "Terraform"
  }

  envs = ["dev", "prod"]

  region = "eastus"

  # USed to make a unique stg account
  stgid = "1542"

  # github repo path (excluding github.com)
  repo = "bobby060/portfolio"

  gh_identities = {
    plan        = "repo:${local.repo}:pull_request"
    deploy-dev  = "repo:${local.repo}:environment:dev"
    deploy-prod = "repo:${local.repo}:environment:prod"
  }


}