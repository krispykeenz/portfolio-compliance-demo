#!/bin/sh
set -eu
base_url="${BASE_URL:-http://localhost:5173}"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

curl --fail --silent --show-error "$base_url/healthz" >/dev/null
printf '%s' '00000000140a12706f7274666f6c696f2d64656d6f2d303031' | xxd -r -p > "$tmp_dir/request.bin"
curl --fail --silent --show-error \
  -H 'Content-Type: application/grpc-web+proto' \
  -H 'X-Grpc-Web: 1' \
  --data-binary "@$tmp_dir/request.bin" \
  "$base_url/compliance.v1.ComplianceService/GetPortfolio" > "$tmp_dir/response.bin"
strings "$tmp_dir/response.bin" | grep -q 'portfolio-demo-001'
echo 'Docker Compose health and gRPC-Web smoke checks passed.'
