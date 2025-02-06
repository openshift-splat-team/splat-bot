#!/bin/bash

set -o nounset
set -o pipefail

export UNIT="TRUE"
REPO_ROOT=$(dirname "${BASH_SOURCE}")/..

echo KUBEBUILDER_ASSETS=$KUBEBUILDER_ASSETS

GINKGO=${GINKGO:-"go run ${REPO_ROOT}/vendor/github.com/onsi/ginkgo/v2/ginkgo"}
GINKGO_ARGS=${GINKGO_ARGS:-"-v --randomize-all --randomize-suites --keep-going --race --trace --timeout=10m"}
GINKGO_EXTRA_ARGS=${GINKGO_EXTRA_ARGS:-""}

# Ensure that some home var is set and that it's not the root.
# This is required for the kubebuilder cache.
export HOME=${HOME:=/tmp/kubebuilder-testing}
if [ $HOME == "/" ]; then
  export HOME=/tmp/kubebuilder-testing
fi

TEST_PACKAGES=${TEST_PACKAGES:-$(go list -f "{{ .Dir }}" ./test/...)}

# Print the command we are going to run as Make would.
echo ${GINKGO} ${GINKGO_ARGS} ${GINKGO_EXTRA_ARGS} "<omitted>"
${GINKGO} ${GINKGO_ARGS} ${GINKGO_EXTRA_ARGS} ${TEST_PACKAGES}
