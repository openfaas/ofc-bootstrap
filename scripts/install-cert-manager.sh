# Add the Jetstack Helm repository
helm repo add jetstack https://charts.jetstack.io

# Update your local Helm chart repository cache
helm repo update

# Install the cert-manager Helm chart
helm upgrade --install \
  cert-manager \
  --namespace cert-manager \
  --version v0.13.0 \
  jetstack/cert-manager \
  --wait 
