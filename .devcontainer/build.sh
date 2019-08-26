#!/usr/bin/env bash
set -euxo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"

source "$ROOT/utils.sh"

cat >> /root/.bashrc <<EOF
if [ -z "\$BUILD_VERSION_CHECK_DONE" ] && ! diff -q "$ROOT" ".devcontainer/" > /dev/null; then
  echo -e "\e[1m\e[31mThis container is outdated. Please rebuild.\e[0m" > /dev/stderr
fi
export BUILD_VERSION_CHECK_DONE=true
EOF

download_tgz https://github.com/goreleaser/goreleaser/releases/download/v0.116.0/goreleaser_Linux_x86_64.tar.gz \
  34b7e3b843158bd0714d1be996951685496491adab4524fb1198ae144ab06ba3 /usr/local/bin \
  goreleaser
