#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

readonly OS_LIST=("darwin" "linux")
readonly ARCH_LIST=("amd64")

if [[ ! -d "./releases" ]]; then
  mkdir releases
fi

for os in "${OS_LIST[@]}"; do
  for arch in "${ARCH_LIST[@]}"; do
    GOOS="${os}" GOARCH="${arch}" go build -o "releases/go-wait-probe-${os}-${arch}"
  done
done

