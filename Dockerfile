ARG GO_VERSION=1.18.3

#######
# UPX #
#######

FROM golang:${GO_VERSION} AS upx

ENV UPX_VERSION=3.96

RUN apt-get update && \
    apt-get install -y xz-utils --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

ADD https://github.com/upx/upx/releases/download/v${UPX_VERSION}/upx-${UPX_VERSION}-amd64_linux.tar.xz /usr/local
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN xz -d -c /usr/local/upx-${UPX_VERSION}-amd64_linux.tar.xz | tar -xOf - upx-${UPX_VERSION}-amd64_linux/upx > /bin/upx && chmod a+x /bin/upx


###########
# MODULES #
###########

FROM golang:${GO_VERSION} AS modules

WORKDIR /src

COPY ./go.mod ./go.sum ./

######## START PRIVATE MODULES with github deploy keys

## reference: https://docs.docker.com/develop/develop-images/build_enhancements/#using-ssh-to-access-private-data-in-builds
## / ! \ you have to create a new deploy key for each repository
## for golang-common, create a new key as well, check other existing keys for naming convention

# Download public key for github.com
RUN mkdir -p -m 0700 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

# Forces the usage of git and ssh key fwded by ssh-agent for monacohq git repos
RUN git config --global url."git@github.com:monacohq/".insteadOf "https://github.com/monacohq/"

######## END PRIVATE MODULES with github deploy keys

# private go packages
ENV GOPRIVATE=github.com/monacohq/*

# --mount=type=ssh only works in buildkit
RUN --mount=type=ssh go mod download


###########
# BUILDER #
###########

FROM golang:${GO_VERSION} AS builder

COPY --from=upx /bin/upx /bin/upx
COPY --from=modules /go/pkg /go/pkg

RUN useradd -u 10001 nonroot

WORKDIR /src

COPY ./ ./

ARG LAST_MAIN_COMMIT_HASH
ARG LAST_MAIN_COMMIT_TIME
ENV FLAG="-X api.CommitTime=${LAST_MAIN_COMMIT_TIME}"
ENV FLAG="$FLAG -X api.CommitHash=${LAST_MAIN_COMMIT_HASH}"

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -ldflags "-s -w $FLAG" \
    -buildvcs=true \
    -o /api ./cmd/api/*.go
    
RUN /bin/upx /api --best --lzma


#########
# FINAL #
#########

FROM scratch AS final

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY ./config /config

COPY --from=builder /api /api

USER nonroot

CMD ["/api"]