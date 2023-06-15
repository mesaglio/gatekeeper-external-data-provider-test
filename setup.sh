#!/bin/bash
helm uninstall external-data-provider --namespace "gatekeeper-system"
# k delete clusterrole gatekeeper-manager-role
# k delete clusterrolebinding gatekeeper-manager-rolebinding
# k delete mutatingwebhookconfigurations gatekeeper-mutating-webhook-configuration
# k delete validatingwebhookconfigurations gatekeeper-validating-webhook-configuration
export NAMESPACE=gatekeeper-system

# generate a self-signed certificate for the external data provider
# ./scripts/generate-tls-cert.sh

# build the image via docker buildx
make docker-build

# load the image into kind
make kind-load-image

helm install external-data-provider charts/external-data-provider \
    --set clientCAFile="" \
    --set provider.tls.caBundle="$(cat certs/ca.crt | base64 | tr -d '\n\r')" \
    --namespace "gatekeeper-system" \
    --create-namespace

kubectl apply -f validation/external-data-provider-constraint-template.yaml && kubectl apply -f validation/external-data-provider-constraint.yaml

# kubectl run nginx --image=error_nginx --dry-run=server -ojson