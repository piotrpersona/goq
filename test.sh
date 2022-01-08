#!/usr/bin/env bash

set -eu

docker build -t goq .
docker run -t goq go build ./...
docker run -t goq go test ./...
