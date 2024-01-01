# Build stage
ARG ALPINE_VERSION
ARG GOLANG_VERSION

FROM ghcr.io/ghcri/golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder
ARG TARGETARCH
ARG TARGETOS
ARG TARGETVARIANT
COPY dist/bazaar_${TARGETOS}_${TARGETARCH}${TARGETVARIANT}/bazaar /usr/bin/bazaar
RUN apk add --no-cache ca-certificates tzdata && \
    chmod +x /usr/bin/bazaar

# Server image
FROM scratch

ENV PORT 8080
LABEL org.opencontainers.image.source="https://github.com/fmartingr/bazaar"
LABEL maintainer="Felipe Martin <github@fmartingr.com>"

COPY --from=builder /usr/bin/bazaar /usr/bin/bazaar
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/usr/bin/bazaar"]
