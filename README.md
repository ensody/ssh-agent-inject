# ssh-agent-inject

Forwards the host's ssh-agent into a Docker container. This is especially useful when working with the [Visual Studio Code Remote - Containers](https://code.visualstudio.com/docs/remote/containers) extension and Git repos clone via SSH.

Why is this needed? While you can bind-mount the `SSH_AUTH_SOCK` from a Linux host, this is [not possible](https://github.com/microsoft/vscode-remote-release/issues/106) from a [macOS](https://github.com/docker/for-mac/issues/410) or Windows host. Also, none of the existing solutions is cross-platform and easy. The [recommended solution](https://code.visualstudio.com/docs/remote/containers#_using-ssh-keys) is to copy the SSH key from the host to the container and enter the password within the container.

This problem is solved by ssh-agent-inject.

## Usage

Add `ENV SSH_AUTH_SOCK=/tmp/.ssh-auth-sock` to your Dockerfile. Label your container with `com.ensody.ssh-agent-inject` (`docker run -l com.ensody.ssh-agent-inject ...`).

Make sure ssh-agent-inject runs in the background or just launch it on-demand.

Note that this project is itself using ssh-agent-inject with VS Code (see `.devcontainer/`).

## Building

All required dependencies are contained in a Docker image defined in `.devcontainer/`, which can be automatically used with Visual Studio Code (or manually via Docker build & run).
That way your host system stays clean and the whole environment is automated, exactly defined, isolated from the host, and easily reproducible.
This saves time and prevents mistakes (wrong version, interference with other software installed on host, etc.).

Use [goreleaser](https://goreleaser.com/) to build binaries for all platforms:

```bash
goreleaser --snapshot --skip-publish --rm-dist
```
