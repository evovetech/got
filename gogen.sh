#!/usr/bin/env bash

pkgs=(
    "collect/map.go"
    "collect/set.go"
    "collect/list.go"
)

for p in "${pkgs[@]}"; do
    go generate -v "./${p}"
done
