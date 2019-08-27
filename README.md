# ssh-agent-inject

[![Build Status](https://travis-ci.com/ensody/ssh-agent-inject.svg?branch=master)](https://travis-ci.com/ensody/ssh-agent-inject)

Forwards the host's ssh-agent into a Docker container. This is especially useful when working with the [Visual Studio Code Remote - Containers](https://code.visualstudio.com/docs/remote/containers) extension and Git repos cloned via SSH.

## Why this is needed

While you can bind-mount the `SSH_AUTH_SOCK` from a Linux host, this is [not possible](https://github.com/microsoft/vscode-remote-release/issues/106) from a [macOS](https://github.com/docker/for-mac/issues/410) or Windows host. Also, none of the existing solutions is cross-platform and easy. The [recommended solution](https://code.visualstudio.com/docs/remote/containers#_using-ssh-keys) is to copy the SSH key from the host to the container, but then you have to manually add the key (assuming you've setup ssh-agent within the container) and enter the password within the container.

With ssh-agent-inject you can skip those annoyances and simply reuse your host's ssh-agent.

## Usage

[Download](https://github.com/ensody/ssh-agent-inject/releases) ssh-agent-inject for your platform. Make sure ssh-agent-inject runs in the background or just launch it on-demand.

Add the following to your Dockerfile:

```Dockerfile
ENV SSH_AUTH_SOCK=/tmp/.ssh-auth-sock
LABEL com.ensody.ssh-agent-inject=
```

Alternatively, you can run an arbitrary container directly:

```
docker run -e SSH_AUTH_SOCK=/tmp/.ssh-auth-sock -l com.ensody.ssh-agent-inject ...
```

Note that this project is itself using ssh-agent-inject with VS Code (see `.devcontainer/`).

## Building

All required dependencies are contained in a Docker image defined in `.devcontainer/`, which can be automatically used with Visual Studio Code (or manually via Docker build & run).
That way your host system stays clean and the whole environment is automated, exactly defined, isolated from the host, and easily reproducible.
This saves time and prevents mistakes (wrong version, interference with other software installed on host, etc.).

Run `./build.sh` to build binaries for all platforms.
