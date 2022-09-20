#!/bin/bash

set -exo pipefail

echo "Creating cluster config"

cat <<EOF > /tmp/kube.config
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ${INPUT_CA_CERT}
    server: ${INPUT_SERVER}
  name: tap
contexts:
- context:
    cluster: tap
    namespace: dev
    user: tap
  name: tap
current-context: tap
kind: Config
preferences: {}
users:
- name: tap
  user:
    token: ${INPUT_TOKEN}
EOF
cat /tmp/kube.config

echo "Authenticating with kubectl"
export KUBECONFIG=/tmp/kube.config 
kubectl get namespaces

echo "Check that the kpack cli is working"
kpack image list

echo "Create the kpack resource and tail the build log"
