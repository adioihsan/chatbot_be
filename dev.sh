#!/usr/bin/env bash
set -euo pipefail

mkdir -p tmp

echo "PATH=$PATH"
echo "GOPATH=$GOPATH"
echo "GOMODCACHE=$GOMODCACHE"
/usr/local/go/bin/go version
air -v || true

# Start Delve (remote attach) in background
dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient --continue ./ &>/dev/null &

# Hot reload (Air will build using absolute /usr/local/go/bin/go from .air.toml)
exec air
