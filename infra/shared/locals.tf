locals {
  state = {
    resource_group_name  = "rgp-portfolio-shared"
    storage_account_name = "stgportfolio1542state"
  }

  # env => hostname. Add prod here when its stack exists.
  hosts = {
    dev = "test.rnorwood.com"
  }
}