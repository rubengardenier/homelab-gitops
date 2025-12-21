# Talos Cluster Provisioning

This guide describes how to install, manage, and recover the Talos Kubernetes cluster using Terraform.

## Prerequisites

- [Terraform](https://www.terraform.io/downloads) >= 1.0
- [talosctl](https://www.talos.dev/latest/introduction/getting-started/#talosctl) CLI
- Talos Linux ISO on USB stick (for initial boot)
- Machines with NVMe disk (`/dev/nvme0n1`)

## File Structure

```
terraform/
├── main.tf                 # Main Terraform configuration
├── variables.tf            # Variable definitions
├── outputs.tf              # Output definitions
├── providers.tf            # Provider configuration
├── terraform.tfvars        # Your IPs (NOT in git!)
├── terraform.tfvars.example # Example configuration
├── patches/
│   ├── controlplane.yaml   # Control plane specific config
│   ├── worker-1.yaml       # Worker 1 config
│   └── worker-2.yaml       # Worker 2 config
└── scripts/
    └── export-configs.sh   # Script to export configs
```

## Initial Installation

### Step 1: Boot machines from Talos ISO

1. Download the [Talos ISO](https://github.com/siderolabs/talos/releases)
2. Write to USB stick with `dd` or Balena Etcher
3. Boot all machines from USB
4. Wait until they are in **maintenance mode** (you'll see a console with IP address)

### Step 2: Configure Terraform

```bash
cd terraform

# Copy example and fill in your IPs
cp terraform.tfvars.example terraform.tfvars
nano terraform.tfvars
```

### Step 3: Run Terraform

```bash
# Initialize Terraform (first time only)
terraform init

# Preview what will happen
terraform plan

# Apply
terraform apply
```

This automatically:
1. Generates cluster secrets (CA, certificates)
2. Pushes configuration to all nodes
3. Nodes reboot and install Talos to disk
4. Bootstraps the cluster
5. Generates kubeconfig

### Step 4: Export configs

```bash
./scripts/export-configs.sh
```

### Step 5: Verify cluster

```bash
# Check Talos nodes
talosctl --talosconfig ./talosconfig health

# Check Kubernetes nodes
KUBECONFIG=./kubeconfig kubectl get nodes
```

## Useful Commands

```bash
# Check Talos version
talosctl --talosconfig ./talosconfig version

# View logs
talosctl --talosconfig ./talosconfig logs kubelet

# Services status
talosctl --talosconfig ./talosconfig services

# Cluster health
talosctl --talosconfig ./talosconfig health

# Kubernetes nodes
KUBECONFIG=./kubeconfig kubectl get nodes -o wide

# Upgrade Talos
talosctl --talosconfig ./talosconfig upgrade --image ghcr.io/siderolabs/installer:v1.11.0
```
