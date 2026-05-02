ARG GO_DEV_DIGEST=sha256:1f11d8ebe336778bc6031a4ab586f7d1fb6c0c33ec15f8c02074bb4f8e32fb6b
FROM mcr.microsoft.com/devcontainers/go@$GO_DEV_DIGEST as utils
RUN go install golang.org/x/pkgsite/cmd/pkgsite@latest

FROM mcr.microsoft.com/devcontainers/go@$GO_DEV_DIGEST
COPY --from=utils /go/bin/* /usr/local/bin/
RUN echo "en_GB.UTF-8 UTF-8" >> /etc/locale.gen && locale-gen
ENV LANG=en_GB.UTF-8
ENV STARSHIP_CONFIG=/etc/starship.toml
COPY assets/starship.toml /etc/starship.toml
COPY assets/starship.sh /usr/local/bin/star.sh
RUN apt update -y && apt upgrade -y
RUN apt install -y npm
RUN npm install -g opencode-ai


