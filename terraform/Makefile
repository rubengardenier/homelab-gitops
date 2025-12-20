.PHONY: init plan apply destroy export-configs health info

# Initialize Terraform
init:
	terraform init

# Plan changes
plan:
	terraform plan

# Apply changes (bootstrap cluster)
apply:
	terraform apply

# Destroy cluster (use with caution!)
destroy:
	terraform destroy

# Export kubeconfig and talosconfig
export-configs:
	./scripts/export-configs.sh

# Check cluster health
health:
	talosctl --talosconfig talosconfig health
	kubectl --kubeconfig kubeconfig get nodes -o wide

# Show cluster info (from terraform outputs)
info:
	@terraform output
