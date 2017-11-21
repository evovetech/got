#!/usr/bin/env bash

echo "$ glide install"
glide install
echo "$ glide brew"
glide brew > brew.lock
