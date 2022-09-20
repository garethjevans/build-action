#!/bin/bash

set -euxo pipefail

echo "Creating cluster config"

cat <<EOF > /tmp/kube.config
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: DATA+OMITTED
    server: https://35.205.192.172
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
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      args: null
      command: gke-gcloud-auth-plugin
      env: null
      interactiveMode: IfAvailable
      provideClusterInfo: false
EOF
cat /tmp/kube.config

echo "Authenticating with kubectl"

echo "Check that the kpack cli is working"

echo "Create the kpack resource and tail the build log"

#    - uses: tale/kubectl-action@v1
#    - run: |
#        echo "${{ secrets.CA_CERT }} " > ca.crt
#        curl -v --cacert ca.crt -H "Authorization: Bearer ${{ secrets.TOKEN }}" ${{ secrets.HOST }}/apis/kpack.io/v1alpha2/namespaces/${{ secrets.NAMESPACE }}/images

