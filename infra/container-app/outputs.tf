output "fqdn" {
  description = "Default ACA hostname (CNAME target for the proxied DNS record)"
  value       = azurerm_container_app.this.ingress[0].fqdn
}

output "custom_domain_verification_id" {
  description = "Value for the asuid.<host> TXT record that ACA uses to verify custom domains"
  value       = azurerm_container_app_environment.this.custom_domain_verification_id
}

output "environment_id" {
  value = azurerm_container_app_environment.this.id
}

output "container_app_id" {
  value = azurerm_container_app.this.id
}
