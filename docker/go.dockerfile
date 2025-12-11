FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy only the `app` directory into the image (avoids building from repo root)
COPY app ./app

# Switch into the app folder where main.go and go.mod live
WORKDIR /app/app/

# If go.mod is missing, initialize it. Then ensure dependencies and build.
# This makes the build tolerant when the repo doesn't include go.mod/go.sum yet.
RUN if [ ! -f go.mod ]; then \
			go mod init github.com/johnkespitia/taller-go-repo/app || true; \
		fi && \
		go env -w GOPROXY=https://proxy.golang.org,direct && \
		go mod tidy && \
		go build -o main ./cmd/

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/app/main .

EXPOSE 8080
CMD ["./main"]

## Dev stage: image for development with `go` and `air` for live reload
FROM golang:1.21-alpine AS dev

ENV PATH=/go/bin:${PATH}

# Install build tools required to `go install` tools like `air`
RUN apk add --no-cache git build-base ca-certificates curl

# Create app dir (we expect host to mount source into this path in dev)
WORKDIR /app/app

# Copy module files only to take advantage of layer caching for dependencies
#COPY app/go.mod app/go.sum ./
RUN if [ -f go.mod ]; then go mod download; fi

# Try to install `air` (live reload). If `air` requires a newer Go, fall back to CompileDaemon.
RUN (go install github.com/cosmtrek/air@latest) || go install github.com/githubnemo/CompileDaemon@latest

# Create a small entrypoint that prefers `air` if available, otherwise runs CompileDaemon.
RUN printf '#!/bin/sh\n\nif command -v air >/dev/null 2>&1; then\n  exec air -c .air.toml\nelse\n  exec CompileDaemon -log-prefix=false -build="go build -o ./tmp/main ./..." -command="./tmp/main"\nfi\n' > /usr/local/bin/dev-entrypoint \
	&& chmod +x /usr/local/bin/dev-entrypoint

# Expose app port and set working dir. In dev you should mount your source
# from host: `./app:/app/app` so changes trigger reloads.
EXPOSE 8080

# Default command runs the dev-entrypoint which selects the available tool
CMD ["/usr/local/bin/dev-entrypoint"]