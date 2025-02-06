# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.29

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
CONTROLLER_GEN = go run ${PROJECT_DIR}/vendor/sigs.k8s.io/controller-tools/cmd/controller-gen
ENVTEST = go run ${PROJECT_DIR}/vendor/sigs.k8s.io/controller-runtime/tools/setup-envtest
GINKGO = go run ${PROJECT_DIR}/vendor/github.com/onsi/ginkgo/v2/ginkgo
GOLANGCI_LINT = go run ${PROJECT_DIR}/vendor/github.com/golangci/golangci-lint/cmd/golangci-lint

VERSION     ?= $(shell git describe --always --abbrev=7)
MUTABLE_TAG ?= latest
IMAGE       ?= cluster-control-plane-machine-set-operator
BUILD_IMAGE ?= registry.ci.openshift.org/openshift/release:golang-1.21

.PHONY: test
test:
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path --bin-dir $(PROJECT_DIR)/bin)" ./hack/test.sh