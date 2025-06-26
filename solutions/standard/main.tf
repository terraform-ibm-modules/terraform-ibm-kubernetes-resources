data "ibm_container_cluster_config" "cluster_config" {
  cluster_name_id   = var.cluster_name
  admin             = true
  resource_group_id = var.resource_group_id
  endpoint_type     = var.cluster_config_endpoint_type != "default" ? var.cluster_config_endpoint_type : null # null value represents default
}


module "namespaces" {
  source     = "terraform-ibm-modules/namespace/ibm"
  version    = "1.0.3"
  namespaces = var.namespaces
}
