#!/bin/bash
set -euo pipefail

mkdir -p test/coverage

# Run unit tests
go test -covermode=set -coverprofile=test/coverage/coverage.out \
  ./internal/biz/usecase/... \
  ./internal/handler/account/... \
  ./internal/handler/jwt/... \
  ./internal/handler/session/... \
  ./internal/handler/health/... 
