#!/bin/bash

set -o pipefail

options=()
[[ -n "$TEST_NAME" ]] && options=(-run "$TEST_NAME")

echo "options = ${options[@]}"

pkgs=()
if [[ -n "$TEST_PKG" ]]; then
  pkgs=("$TEST_PKG")
else
  pkgs=('./...')
fi

go test -v "${pkgs[@]}" -bench . "${options[@]}" | awk -f ./scripts/go_test_color_status.awk
