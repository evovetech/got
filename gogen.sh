#!/usr/bin/env bash

pkgs=(
    "cmd/merge/mv/file/node.go"
    "cmd/merge/mv/file/nodeList.go"
)

for p in "${pkgs[@]}"; do
    go generate -v "./${p}"
done
