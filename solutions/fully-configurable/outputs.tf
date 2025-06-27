data "kubernetes_all_namespaces" "allns" {
  depends_on = [module.namespaces]
}

output "all_namespaces" {
  value       = data.kubernetes_all_namespaces.allns.namespaces
  description = "Returns all namespaces"
}
