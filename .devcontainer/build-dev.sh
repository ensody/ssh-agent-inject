#!/usr/bin/env bash
set -euxo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"

"$ROOT/build-common.sh"

source "$ROOT/utils.sh"

cat >> ~/.bashrc <<EOF
if [ -z "\$BUILD_VERSION_CHECK_DONE" ] && ! diff -q "$ROOT" ".devcontainer/" > /dev/null; then
  echo -e "\e[1m\e[31mThis container is outdated. Please rebuild.\e[0m" > /dev/stderr
fi
export BUILD_VERSION_CHECK_DONE=true
EOF

# Make sure a few common tools are installed
apt-get update
apt-get upgrade -y
apt-get install -y --no-install-recommends curl gettext git gnupg less procps apt-utils locales

# This fixed the locale and is required for diff-so-fancy
echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen
locale-gen en_US.UTF-8

# Install diff-so-fancy
DIFF_SO_FANCY_VERSION="1.2.6"
download https://raw.githubusercontent.com/so-fancy/diff-so-fancy/v${DIFF_SO_FANCY_VERSION}/third_party/build_fatpack/diff-so-fancy \
  ed9de2669c789d1aba8456d0a7cf95adb326e220c99af4336405f21add8f0852 /usr/bin/diff-so-fancy
chmod a+x /usr/bin/diff-so-fancy

# The following is taken from
# https://github.com/microsoft/vscode-remote-try-go/blob/master/.devcontainer/Dockerfile

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
