#!/usr/bin/env bash

pkgs=(
    "collect/set.go"
)

for p in "${pkgs[@]}"; do
    go generate -v "./${p}"
done
