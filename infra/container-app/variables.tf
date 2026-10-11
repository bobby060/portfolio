variable "resource_group_name" {
    description = "Name of resource group in which to create container app"
    type = string
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

variable "container" {
    type = string
    description = "image path to container to run ex ghcr.io/image"
}