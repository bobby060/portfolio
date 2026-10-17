variable "cloudflare_zone_id" {
  description = "Zone ID for rnorwood.com"
  type        = string
}

variable "cloudflare_account_id" {
  description = "Cloudflare account ID (needed for Zero Trust / Access)"
  type        = string
}

variable "custom_domain" {
  description = "Hostname for the dev app"
  type        = string
  default     = "dev.rnorwood.com"
}

variable "container_path" {
  description = "Image to run, e.g. ghcr.io/bobby060/site@sha256:..."
  type        = string
}

variable "access_allowed_emails" {
  description = "Emails allowed through Cloudflare Access. Set in terraform.tfvars or TF_VAR_access_allowed_emails."
  type        = list(string)
}

variable "dns_proxied" {
  description = <<-EOT
    Proxy the dev hostname through Cloudflare. Start false for the first apply
    so ACA can validate the CNAME, then set true (DESIGN.md §3.2).
  EOT
  type        = bool
  default     = false
}
