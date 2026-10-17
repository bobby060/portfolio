variable "resource_group_name" {
  description = "Name of resource group in which to create container app"
  type        = string
}

variable "name" {
  type = string
}

variable "tags" {
  type = map(string)
}

variable "containername" {
  type = string
}

variable "containerpath" {
  type        = string
  description = "image path to container to run ex ghcr.io/image"
}

variable "target_port" {
  type        = number
  description = "Port the container listens on (site/Dockerfile EXPOSE 8080)"
  default     = 8080
}

variable "ingress_allowed_cidrs" {
  type        = list(string)
  description = "CIDRs allowed to reach ingress. Pass Cloudflare's ranges so only the proxy can connect."
}
