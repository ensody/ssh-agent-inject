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

# Install gocode-gomod
go get -x -d github.com/stamblerre/gocode 2>&1
go build -o gocode-gomod github.com/stamblerre/gocode
mv gocode-gomod $GOPATH/bin/

# Install Go tools
go get -u -v \
  github.com/mdempsky/gocode \
  github.com/uudashr/gopkgs/cmd/gopkgs \
  github.com/ramya-rao-a/go-outline \
  github.com/acroca/go-symbols \
  github.com/godoctor/godoctor \
  golang.org/x/tools/cmd/guru \
  golang.org/x/tools/cmd/gorename \
  github.com/rogpeppe/godef \
  github.com/zmb3/gogetdoc \
  github.com/haya14busa/goplay/cmd/goplay \
  github.com/sqs/goreturns \
  github.com/josharian/impl \
  github.com/davidrjenni/reftools/cmd/fillstruct \
  github.com/fatih/gomodifytags \
  github.com/cweill/gotests/... \
  golang.org/x/tools/cmd/goimports \
  golang.org/x/lint/golint \
  golang.org/x/tools/cmd/gopls \
  github.com/alecthomas/gometalinter \
  honnef.co/go/tools/... \
  github.com/golangci/golangci-lint/cmd/golangci-lint \
  github.com/mgechev/revive \
  github.com/derekparker/delve/cmd/dlv

download_tgz https://github.com/goreleaser/goreleaser/releases/download/v0.116.0/goreleaser_Linux_x86_64.tar.gz \
  34b7e3b843158bd0714d1be996951685496491adab4524fb1198ae144ab06ba3 /usr/local/bin \
  goreleaser
