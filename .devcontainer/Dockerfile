ARG GO_VERSION=dev-1.25
FROM mcr.microsoft.com/devcontainers/go:$GO_VERSION as utils
RUN go install golang.org/x/pkgsite/cmd/pkgsite@latest

FROM mcr.microsoft.com/devcontainers/go:$GO_VERSION
COPY --from=utils /go/bin/* /usr/local/bin/
RUN echo "en_GB.UTF-8 UTF-8" >> /etc/locale.gen && locale-gen
ENV LANG=en_GB.UTF-8
ENV STARSHIP_CONFIG=/etc/starship.toml
COPY assets/starship.toml /etc/starship.toml
COPY assets/starship.sh /usr/local/bin/star.sh
RUN apt update -y && apt upgrade -y
RUN apt install -y npm
RUN npm install -g opencode-ai


