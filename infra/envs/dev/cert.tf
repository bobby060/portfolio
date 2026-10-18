# ---- Origin TLS: Cloudflare Origin CA cert, uploaded to the ACA environment ----
resource "tls_private_key" "origin" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P256"
}

resource "tls_cert_request" "origin" {
  private_key_pem = tls_private_key.origin.private_key_pem

  subject {
    common_name = var.custom_domain
  }
}

resource "cloudflare_origin_ca_certificate" "origin" {
  csr                = tls_cert_request.origin.cert_request_pem
  hostnames          = [var.custom_domain]
  request_type       = "origin-ecc"
  requested_validity = 5475 # 15 years
}

# PEM bundle (cert + PKCS#8 key); PEM uploads must not have a password.
resource "azurerm_container_app_environment_certificate" "origin" {
  name                         = "${local.name}-origin-ca"
  container_app_environment_id = module.container_app.environment_id
  certificate_blob_base64      = base64encode("${trimspace(cloudflare_origin_ca_certificate.origin.certificate)}\n${tls_private_key.origin.private_key_pem_pkcs8}")
  certificate_password         = ""
}
