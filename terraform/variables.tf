variable "cluster_name" {
  description = "Name of the Talos cluster"
  type        = string
  default     = "talos-cluster"
}

variable "cluster_endpoint" {
  description = "Kubernetes API endpoint (e.g., https://10.0.0.10:6443)"
  type        = string
}

variable "controlplane_nodes" {
  description = "Control plane node IPs"
  type        = list(string)
}

variable "worker_nodes" {
  description = "Worker node IPs"
  type        = list(string)
}

variable "talos_version" {
  description = "Talos version to use"
  type        = string
  default     = "v1.11.0"
}

variable "kubernetes_version" {
  description = "Kubernetes version to use"
  type        = string
  default     = "1.32.1"
}
