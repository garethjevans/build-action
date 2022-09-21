#!/bin/bash

set -euo pipefail

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
    namespace: ${INPUT_NAMESPACE}
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
kubectl get images.kpack.io -n $INPUT_NAMESPACE

echo "Check that the kpack cli is working"
kp image list

echo "Create the kpack resource and tail the build log"
IMAGE_NAME=$(echo $GITHUB_REPOSITORY | sed 's|/|-|')
kp images create $IMAGE_NAME \
	--tag $INPUT_DESTINATION \
	--git $GITHUB_SERVER_URL/$GITHUB_REPOSITORY \
	--git-revision $GITHUB_SHA \
	--namespace $INPUT_NAMESPACE

echo "Image $IMAGE_NAME Created"
# TODO --wait does not seem to work, possibly due to a lack of tty

trap "kubectl delete images.kpack.io $IMAGE_NAME" EXIT

counter=0

until [ $counter -gt 100 ]
do
  set +e
  STATUS=$(kubectl get images.kpack.io $IMAGE_NAME --namespace $INPUT_NAMESPACE -ojsonpath="{.status.conditions[?(@.type=='Ready')].status}")
  echo "Check $counter> $STATUS"
  set -e

  if [[ "$STATUS" == "True" ]]; then
    break
  fi
  
  if [[ "$STATUS" == "False" ]]; then
    echo "$IMAGE_NAME failed to become ready"
    exit 1
  fi

  set +e
  ((counter++))
  set -e 
  sleep 5
done

# TODO how do we determine if this has passed / failed
# What happens with a timeout?

BUILT_IMAGE_NAME=$(kubectl get images.kpack.io $IMAGE_NAME -ojsonpath="{.status.latestImage}")
echo "::set-output name=name::$BUILT_IMAGE_NAME"



