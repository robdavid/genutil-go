ARG GO_DEV_DIGEST=sha256:0027553e7c9696d975bc28afdb254e8c4799fe0df62c75565d33684c4ed602e4
# Version dev-1.26
FROM mcr.microsoft.com/devcontainers/go@$GO_DEV_DIGEST AS utils
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


