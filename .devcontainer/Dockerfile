FROM golang:1.12

ENV DEBIAN_FRONTEND=noninteractive

# We use a fixed path for ssh-agent injection and the Git UI of Visual Studio Code
ENV SSH_AUTH_SOCK /tmp/.ssh-auth-sock

COPY . /dockerbuild/ssh-agent-inject
RUN /dockerbuild/ssh-agent-inject/build.sh
