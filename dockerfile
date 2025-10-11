
FROM golang:1.22.3 AS builder

WORKDIR /src

# Update OS packages and install minimal build deps + CA bundle
COPY go.mod ./
RUN apt-get update && \
	apt-get upgrade -y && \
	apt-get install -y --no-install-recommends ca-certificates git && \
	update-ca-certificates && \
	rm -rf /var/lib/apt/lists/* && \
	go mod download || true

# copy source and build a static, stripped binary (smaller, no cgo)
COPY . .
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -trimpath -ldflags="-s -w" -o /src/main ./cmd

FROM gcr.io/distroless/static:nonroot

# Copy binary and CA certs from builder. Distroless nonroot has a minimal filesystem and runs as non-root.
COPY --from=builder /src/main /app/main
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["/app/main"]
