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

## Cluster Recovery / Reinstallation

If you want to reinstall the cluster:

### Option A: With working talosconfig (remote reset)

If you still have a working `talosconfig` that matches the nodes:

```bash
# Reset all nodes to maintenance mode
talosctl --talosconfig ./talosconfig reset --graceful=false --reboot --nodes <NODE_IP>

# Wait until they're in maintenance mode, then:
rm terraform.tfstate terraform.tfstate.backup talosconfig kubeconfig
terraform apply
./scripts/export-configs.sh
```

### Option B: Physical access required

If certificates don't match (e.g., after losing terraform state):

1. **Boot all machines from Talos USB**
   - Select boot from USB in BIOS (F12/F8/Esc)
   - If Talos says there's already an installation: format the disk first

2. **Wipe disk (if needed)**
   - Boot from Ubuntu Live USB or similar
   - `sudo wipefs -a /dev/nvme0n1`
   - Reboot from Talos USB

3. **Remove old state**
   ```bash
   rm -f terraform.tfstate terraform.tfstate.backup talosconfig kubeconfig
   ```

4. **Terraform apply**
   ```bash
   terraform apply
   ./scripts/export-configs.sh
   ```

## Adding a New Worker

1. **Create patch file**
   ```bash
   cp patches/worker-2.yaml patches/worker-3.yaml
   # Change hostname to talos-worker-3
   ```

2. **Add IP to terraform.tfvars**

3. **Boot new machine from Talos ISO**

4. **Apply**
   ```bash
   terraform apply
   ```

## Important Files (DO NOT commit!)

These files are in `.gitignore` and contain sensitive data:

| File | Contents |
|------|----------|
| `terraform.tfstate` | Full cluster state including secrets |
| `terraform.tfvars` | IP addresses |
| `talosconfig` | Talos API certificates |
| `kubeconfig` | Kubernetes API certificates |

**Tip:** Back up `terraform.tfstate` and `talosconfig` to a secure location!

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
