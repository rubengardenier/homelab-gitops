# Talos Cluster Secrets
resource "talos_machine_secrets" "this" {}

# Generate client configuration
data "talos_client_configuration" "this" {
  cluster_name         = var.cluster_name
  client_configuration = talos_machine_secrets.this.client_configuration
  endpoints            = var.controlplane_nodes
  nodes                = concat(var.controlplane_nodes, var.worker_nodes)
}

# Control Plane Machine Configuration
data "talos_machine_configuration" "controlplane" {
  cluster_name     = var.cluster_name
  cluster_endpoint = var.cluster_endpoint
  machine_type     = "controlplane"
  machine_secrets  = talos_machine_secrets.this.machine_secrets

  talos_version      = var.talos_version
  kubernetes_version = var.kubernetes_version

  config_patches = [
    file("${path.module}/patches/controlplane.yaml")
  ]
}

# Worker Machine Configurations (one per worker for unique hostnames)
data "talos_machine_configuration" "worker" {
  for_each = { for idx, ip in var.worker_nodes : ip => idx + 1 }

  cluster_name     = var.cluster_name
  cluster_endpoint = var.cluster_endpoint
  machine_type     = "worker"
  machine_secrets  = talos_machine_secrets.this.machine_secrets

  talos_version      = var.talos_version
  kubernetes_version = var.kubernetes_version

  config_patches = [
    file("${path.module}/patches/worker-${each.value}.yaml")
  ]
}

# Apply configuration to control plane nodes
resource "talos_machine_configuration_apply" "controlplane" {
  for_each = toset(var.controlplane_nodes)

  client_configuration        = talos_machine_secrets.this.client_configuration
  machine_configuration_input = data.talos_machine_configuration.controlplane.machine_configuration
  node                        = each.value
  endpoint                    = each.value

  # Use insecure mode for initial apply to nodes in maintenance mode
  apply_mode = "reboot"
}

# Apply configuration to worker nodes
resource "talos_machine_configuration_apply" "worker" {
  for_each = { for idx, ip in var.worker_nodes : ip => idx + 1 }

  client_configuration        = talos_machine_secrets.this.client_configuration
  machine_configuration_input = data.talos_machine_configuration.worker[each.key].machine_configuration
  node                        = each.key
  endpoint                    = each.key

  # Use insecure mode for initial apply to nodes in maintenance mode
  apply_mode = "reboot"
}

# Bootstrap the cluster (only needed once)
resource "talos_machine_bootstrap" "this" {
  depends_on = [talos_machine_configuration_apply.controlplane]

  client_configuration = talos_machine_secrets.this.client_configuration
  node                 = var.controlplane_nodes[0]
  endpoint             = var.controlplane_nodes[0]
}

# Get kubeconfig
data "talos_cluster_kubeconfig" "this" {
  depends_on = [talos_machine_bootstrap.this]

  client_configuration = talos_machine_secrets.this.client_configuration
  node                 = var.controlplane_nodes[0]
  endpoint             = var.controlplane_nodes[0]
}
