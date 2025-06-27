#######################################################################################################################
# Cluster Config
#######################################################################################################################
data "ibm_container_cluster_config" "cluster_config" {
  cluster_name_id = var.cluster_name
  endpoint_type   = var.cluster_config_endpoint_type != "default" ? var.cluster_config_endpoint_type : null # null value represents default
}

#######################################################################################################################
# Namespaces
#######################################################################################################################

module "namespaces" {
  source     = "terraform-ibm-modules/namespace/ibm"
  version    = "1.0.3"
  namespaces = var.namespaces
}
