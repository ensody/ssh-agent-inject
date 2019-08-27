#!/usr/bin/env bash
set -euxo pipefail

export SSH_AGENT_PID=
cleanup() {
  set +e
  docker kill ssh-agent-inject-test
  kill -9 $SSH_AGENT_PID $(jobs -p)
}
trap cleanup INT TERM EXIT

docker run -d --name ssh-agent-inject-test --rm \
  -e SSH_AUTH_SOCK=/tmp/.ssh-auth-inject -l com.ensody.ssh-agent-inject \
  alpine sh -c 'apk add --no-cache openssh-client && sleep 1800'

# Test in a clean environment with custom ssh-agent and ssh-key
eval "$(ssh-agent)"
tmpkey="$(mktemp)"
! yes y | ssh-keygen -t ed25519 -N "" -C "$tmpkey" -f "$tmpkey"
ssh-add "$tmpkey"
rm "$tmpkey"

# Inject ssh-agent and wait for injection to finish
./dist/unix_linux_amd64/ssh-agent-inject -v &
sleep 2

docker exec -t ssh-agent-inject-test ssh-add -l | grep -i ed25519 | grep "$tmpkey"
