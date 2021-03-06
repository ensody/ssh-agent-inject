#!/usr/bin/env bash
set -euxo pipefail

DIR="$(cd "$(dirname "$0")" && pwd)"

source "$DIR/utils.sh"

mkdir -p ~/bin

download_tgz https://github.com/goreleaser/goreleaser/releases/download/v0.116.0/goreleaser_Linux_x86_64.tar.gz \
  34b7e3b843158bd0714d1be996951685496491adab4524fb1198ae144ab06ba3 ~/bin \
  goreleaser
