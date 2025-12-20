output "talosconfig" {
  description = "Talos client configuration"
  value       = data.talos_client_configuration.this.talos_config
  sensitive   = true
}

output "kubeconfig" {
  description = "Kubernetes kubeconfig"
  value       = data.talos_cluster_kubeconfig.this.kubeconfig_raw
  sensitive   = true
}

output "cluster_endpoint" {
  description = "Kubernetes API endpoint"
  value       = var.cluster_endpoint
}

output "controlplane_nodes" {
  description = "Control plane node IPs"
  value       = var.controlplane_nodes
}

output "worker_nodes" {
  description = "Worker node IPs"
  value       = var.worker_nodes
}
