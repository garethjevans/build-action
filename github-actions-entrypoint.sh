#!/bin/bash

set -euo pipefail

echo "Creating cluster config"

CA_CERT_BASE64=$(echo -n "${INPUT_CA_CERT}" | base64 -w 0)
cat <<EOF > /tmp/kube.config
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: ${CA_CERT_BASE64}
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

echo "Authenticating with kubectl"
export KUBECONFIG=/tmp/kube.config 
kubectl get images.kpack.io -n dev
kubectl config view --minify

echo "Check that the kpack cli is working"
kp image list

echo "Create the kpack resource and tail the build log"
IMAGE_NAME=$(echo $GITHUB_REPOSITORY | sed 's|/|-|')
kp images create $IMAGE_NAME \
	--tag $INPUT_DESTINATION \
	--git $GITHUB_SERVER_URL/$GITHUB_REPOSITORY \
	--git-revision $GITHUB_SHA \
	--wait

# TODO how do we determine if this has passed / failed

BUILT_IMAGE_NAME=$(kubectl get images.kpack.io $IMAGE_NAME -ojsonpath="{.status.latestImage}")
echo '::set-output name=name::$BUILT_IMAGE_NAME'

kubectl delete images.kpack.io $IMAGE_NAME


