#!/usr/bin/env bash

paths=(
    "cmd"
    "git"
    "got"
    "util"
    "main.go"
)

function gofmt() {
    echo "$ gofmt $( pwd )/$1"
    command gofmt -w -s "$1"
}

function gofmt_paths() {
    local f
    for f in "$@"; do
        gofmt "${f}"
    done
    return 0
}

gofmt_paths "${paths[@]}"
