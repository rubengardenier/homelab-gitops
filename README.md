# homelab-gitops

GitOps repository for my k3s homelab, managed with FluxCD.

## Repository structure

- `clusters/` – Flux bootstrap and per-cluster Kustomizations (currently `staging` with a future `production` placeholder).
- `apps/` – Application manifests (e.g. Linkding, Mealie) with base and environment overlays.
- `infrastructure/` – Shared platform services, such as networking and Cloudflare Tunnel.
- `monitoring/` – kube-prometheus-stack (Prometheus, Alertmanager, Grafana) plus cluster-specific configuration.

## Tech stack & focus

- k3s, FluxCD, Kustomize, Helm
- GitOps-driven deployments and upgrades (incl. Renovate)
- SOPS-encrypted secrets and Cloudflare-backed ingress
- Goal: treat the homelab as production-like IaC to practice cloud-native platform engineering.
