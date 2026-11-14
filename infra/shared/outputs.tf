output "origin_auth" {
  value       = random_password.origin_auth
  description = "Map of origin auth values from cloudflare"
  sensitive   = true
}