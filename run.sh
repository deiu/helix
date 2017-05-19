#!/bin/sh

export HELIX_PORT="3000"
export HELIX_HOST="localhost"
export HELIX_ROOT="."
export HELIX_STATIC_DIR="./static"
export HELIX_LOG=""
export HELIX_DEBUG=""
export HELIX_CERT="test_cert.pem"
export HELIX_KEY="test_key.pem"

go run bin/*.go
