FROM mcr.microsoft.com/devcontainers/go
RUN echo "en_GB.UTF-8 UTF-8" >> /etc/locale.gen && locale-gen
ENV LANG=en_GB.UTF-8

USER vscode
RUN go install golang.org/x/tools/cmd/godoc@latest
RUN go install golang.org/x/pkgsite/cmd/pkgsite@latest
COPY assets/bashrc /home/vscode/.bashrc

USER root
ENV STARSHIP_CONFIG=/etc/starship.toml
COPY assets/starship.toml /etc/starship.toml
