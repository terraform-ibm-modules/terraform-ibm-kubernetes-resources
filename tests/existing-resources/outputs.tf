output "cluster_name" {
  description = "Name of the cluster"
  value       = ibm_container_vpc_cluster.cluster.name
}

output "resource_group_id" {

  description = "ID of the resource group in which OCP cluster was deployed"
  value       = module.resource_group.resource_group_id
}
