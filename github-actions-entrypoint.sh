#!/usr/bin/env bash

set -euo pipefail

export CA_CERT=${INPUT_CA_CERT}
export SERVER=${INPUT_SERVER}
export NAMESPACE=${INPUT_NAMESPACE}
export TOKEN=${INPUT_TOKEN}
export TAG=${INPUT_DESTINATION}

/usr/bin/builder

