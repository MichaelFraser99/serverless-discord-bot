variable "name" {
  type = string
  description = "The name pre-pended to all resources in this deployment"
}

variable "stage_name" {
  type = string
  description = "The name of the api gateway stage (e.g. dev, staging, prod)"
  default = "v1"
}

variable "public_key" {
  type = string
  description = "The discord bots public API key"
}

variable "application_id" {
  type = string
  description = "The discord bots application id"
}