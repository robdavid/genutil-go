ARG GO_DEV_DIGEST=1.27@sha256:8f1c83128e4bce9d745b1dc0f28f8c027e303d3886e2434966ff135296c41770
FROM mcr.microsoft.com/devcontainers/go:$GO_DEV_DIGEST as utils
RUN go install golang.org/x/pkgsite/cmd/pkgsite@latest

FROM mcr.microsoft.com/devcontainers/go:$GO_DEV_DIGEST
COPY --from=utils /go/bin/* /usr/local/bin/
RUN echo "en_GB.UTF-8 UTF-8" >> /etc/locale.gen && locale-gen
ENV LANG=en_GB.UTF-8
ENV STARSHIP_CONFIG=/etc/starship.toml
COPY assets/starship.toml /etc/starship.toml
COPY assets/starship.sh /usr/local/bin/star.sh
RUN apt update -y && apt upgrade -y
RUN apt install -y npm
RUN npm install -g opencode-ai


