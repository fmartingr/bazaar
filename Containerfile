# Build stage
ARG GOLANG_VERSION
ARG ALPINE_VERSION

FROM docker.io/library/golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder
WORKDIR /src
COPY . .
RUN apk add --no-cache ca-certificates tzdata make && \
   make build

# Server image
FROM scratch

ENV PORT 8080
LABEL org.opencontainers.image.source="https://github.com/fmartingr/bazaar"
LABEL maintainer="Felipe Martin <github@fmartingr.com>"

COPY --from=builder /src/build/bazaar /usr/bin/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/usr/bin/bazaar"]
