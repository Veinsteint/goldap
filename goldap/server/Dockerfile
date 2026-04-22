# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /build
RUN apk add --no-cache git make ca-certificates tzdata
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
ARG BUILD_TIME=unknown
ARG GIT_COMMIT=unknown
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o goldap-server .

# Runtime stage
FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -g '' appuser && \
    mkdir -p /app/data/logs && \
    chown -R appuser:appuser /app
COPY --from=builder /build/goldap-server /app/
COPY --from=builder /build/config.yml /app/
USER appuser
EXPOSE 8150
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -q --spider http://localhost:8150/api/base/ping || exit 1
ENTRYPOINT ["./goldap-server"]
