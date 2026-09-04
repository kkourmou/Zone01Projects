# syntax=docker/dockerfile:1

ARG GO_VERSION=1.24

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src

ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0

COPY go.mod ./
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/ascii-art-web .

FROM alpine:3.21
WORKDIR /app

# Install bash (alpine includes sh by default via busybox)
RUN apk add --no-cache bash

ARG VCS_REF=""
ARG BUILD_DATE=""
ARG SOURCE_URL=""

LABEL org.opencontainers.image.title="ascii-art-web-dockerize" \
    org.opencontainers.image.description="Go (stdlib) web server that renders ASCII art banners" \
    org.opencontainers.image.source=$SOURCE_URL \
    org.opencontainers.image.revision=$VCS_REF \
    org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.licenses="UNLICENSED"

COPY --from=builder /out/ascii-art-web ./ascii-art-web
COPY --from=builder /src/static ./static
COPY --from=builder /src/templates ./templates
COPY --from=builder /src/ascii-style ./ascii-style

EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["./ascii-art-web"]
