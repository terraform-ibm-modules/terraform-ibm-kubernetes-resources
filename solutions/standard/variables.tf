variable "ibmcloud_api_key" {
  type        = string
  description = "The IBM Cloud api key."
  sensitive   = true
}

variable "region" {
  type        = string
  description = "The region of the existing cluster in which you want to create resources."
  default     = "us-south"
}

variable "provider_visibility" {
  description = "Set the visibility value for the IBM terraform provider. Supported values are `public`, `private`, `public-and-private`. [Learn more](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/guides/custom-service-endpoints)."
  type        = string
  default     = "private"
  validation {
    condition     = contains(["public", "private", "public-and-private"], var.provider_visibility)
    error_message = "Invalid visibility option. Allowed values are `public`, `private`, or `public-and-private`."
  }
}

variable "cluster_name" {
  type        = string
  description = "The name of the existing cluster in which you want to create resources."
}

variable "resource_group_id" {
  type        = string
  description = "The ID of the resource group in which cluster exists."
  default     = ""
}

variable "cluster_config_endpoint_type" {
  description = "Specify which type of endpoint to use for cluster config access: 'default', 'private', 'vpe', 'link'. A 'default' value uses the default endpoint of the cluster."
  type        = string
  default     = "default"
  nullable    = false
}

variable "namespaces" {
  type = list(object({
    name = string
    metadata = optional(object({
      labels      = map(string)
      annotations = map(string)
    }))
  }))
  description = "Set of namespaces to create. [Learn more](https://github.com/terraform-ibm-modules/terraform-ibm-kubernetes-resources/blob/main/solutions/standard/DA-types.md#namespaces-)"
}
