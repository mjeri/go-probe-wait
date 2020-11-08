#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

readonly OS_LIST=("darwin" "linux")
readonly ARCH_LIST=("amd64")

for os in "${OS_LIST[@]}"; do
  for arch in "${ARCH_LIST[@]}"; do
    GOOS="${os}" GOARCH="${arch}" go build -o "bin/${arch}/${os}/go-wait-for-it"
  done
done

