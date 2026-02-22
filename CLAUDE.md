# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

GitOps repository for a Talos Kubernetes homelab managed with FluxCD. The cluster runs on Talos Linux with Terraform-based provisioning.

## Common Commands

### Terraform (cluster provisioning)
```bash
cd terraform
make init          # Initialize Terraform
make plan          # Preview changes
make apply         # Apply changes (bootstrap/update cluster)
make export-configs # Export kubeconfig and talosconfig
make health        # Check cluster health
make info          # Show Terraform outputs
```

### Talos cluster management
```bash
talosctl --talosconfig terraform/talosconfig health
talosctl --talosconfig terraform/talosconfig version
talosctl --talosconfig terraform/talosconfig services
talosctl --talosconfig terraform/talosconfig logs kubelet
```

### Kubernetes
```bash
KUBECONFIG=terraform/kubeconfig kubectl get nodes -o wide
```

### Flux GitOps
```bash
flux reconcile kustomization flux-system --with-source
flux get kustomizations
flux logs
```

## Architecture

### Directory Structure
- `terraform/` - Talos cluster provisioning (IPs and secrets are gitignored)
- `clusters/staging/` - Flux bootstrap and Kustomization entry points
- `apps/` - Application manifests with base/staging overlays
- `infrastructure/` - Platform services (CSI drivers, Cloudflare Tunnel, Renovate)
- `monitoring/` - kube-prometheus-stack and metrics-server

### GitOps Flow
1. `clusters/staging/flux-system/` bootstraps Flux
2. Flux watches `clusters/staging/*.yaml` which reference paths in other directories
3. Dependencies: `infrastructure-controllers` → `networking` → `apps`
4. SOPS with age encryption for secrets (`.sops.yaml` files define encryption rules)

### Kustomize Pattern
Each component uses base/staging overlay structure:
- `apps/base/<app>/` - Base manifests (namespace, deployment, service, storage)
- `apps/staging/<app>/kustomization.yaml` - Environment-specific overrides and secrets

### Terraform Talos Setup
- Uses `siderolabs/talos` provider
- Per-node config patches in `terraform/patches/` (controlplane-N.yaml, worker-N.yaml)
- Generates secrets, pushes config, bootstraps cluster, outputs kubeconfig

## Key Technologies
- Talos Linux (immutable Kubernetes OS)
- FluxCD (GitOps operator)
- Kustomize (manifest management)
- Helm via Flux HelmReleases
- SOPS + age (secret encryption)
- NFS CSI driver for persistent storage
- Cloudflare Tunnel for ingress
