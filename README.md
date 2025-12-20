# homelab-gitops

GitOps repository for my Talos Kubernetes homelab, managed with FluxCD.

## Repository structure

- `terraform/` – Talos cluster provisioning with Terraform (IPs and secrets are gitignored).
- `clusters/` – Flux bootstrap and per-cluster Kustomizations (currently `staging` with a future `production` placeholder).
- `apps/` – Application manifests (e.g. Linkding, Mealie) with base and environment overlays.
- `infrastructure/` – Shared platform services, such as networking and Cloudflare Tunnel.
- `monitoring/` – kube-prometheus-stack (Prometheus, Alertmanager, Grafana) plus cluster-specific configuration.

## Cluster provisioning

The Talos cluster is provisioned using Terraform. Sensitive files (IPs, state, configs) are gitignored.

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars  # Fill in your IPs
make init
make apply
make export-configs
```

## Tech stack & focus

- Talos Linux, Kubernetes, FluxCD, Kustomize, Helm
- Terraform for cluster provisioning
- GitOps-driven deployments and upgrades (incl. Renovate)
- SOPS-encrypted secrets and Cloudflare-backed ingress
- Goal: treat the homelab as production-like IaC to practice cloud-native platform engineering.
