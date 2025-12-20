#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

echo "Exporting talosconfig..."
terraform output -raw talosconfig > talosconfig
chmod 600 talosconfig

echo "Exporting kubeconfig..."
terraform output -raw kubeconfig > kubeconfig
chmod 600 kubeconfig

echo ""
echo "Configs exported successfully!"
echo ""
echo "Usage:"
echo "  export TALOSCONFIG=$PROJECT_DIR/talosconfig"
echo "  export KUBECONFIG=$PROJECT_DIR/kubeconfig"
echo ""
echo "Or add to your shell profile:"
echo "  talosctl --talosconfig $PROJECT_DIR/talosconfig <command>"
echo "  kubectl --kubeconfig $PROJECT_DIR/kubeconfig <command>"
