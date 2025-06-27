data "kubernetes_all_namespaces" "allns" {
  depends_on = [module.namespaces]
}

output "my_namespace_present" {
  value       = data.kubernetes_all_namespaces.allns.namespaces
  description = "Returns all namespaces"
}
