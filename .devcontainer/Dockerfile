FROM golang:1.12.9

ENV DEBIAN_FRONTEND=noninteractive

ENV PATH=$PATH:~/bin

# We use a fixed path for ssh-agent injection and the Git UI of Visual Studio Code
ENV SSH_AUTH_SOCK=/tmp/.ssh-auth-sock
# Required by ssh-agent-inject
LABEL com.ensody.ssh-agent-inject=

COPY . /dockerbuild/ssh-agent-inject
RUN /dockerbuild/ssh-agent-inject/build-dev.sh
