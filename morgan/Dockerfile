# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy shared dependencies
COPY framework/ ./framework/
COPY libraries/ ./libraries/

# Copy service dependencies
COPY morgan/go.mod morgan/go.sum ./morgan/

WORKDIR /app/morgan
RUN go mod download

# Copy service code
COPY morgan/ .

# Build
ARG VERSION=dev
ARG COMMIT=none
ARG BRANCH=unknown
ARG BUILD_TIME=unknown

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X 'github.com/siakup/morgan-be/version.Version=${VERSION}' \
              -X 'github.com/siakup/morgan-be/version.Commit=${COMMIT}' \
              -X 'github.com/siakup/morgan-be/version.Branch=${BRANCH}' \
              -X 'github.com/siakup/morgan-be/version.BuildTime=${BUILD_TIME}'" \
    -o /go-morgan main.go

# Stage 2: Development (with air for live reload)
FROM golang:1.25-alpine AS dev

WORKDIR /app

# Install air for live reload (updated package path)
RUN go install github.com/air-verse/air@latest

# Copy shared dependencies
COPY framework/ ./framework/
COPY libraries/ ./libraries/

# Copy service (code will be mounted via volume in docker-compose)
COPY go.work ./
COPY morgan/go.mod morgan/go.sum ./morgan/
COPY morgan/ ./morgan/
COPY morgan/.air.toml ./morgan/

WORKDIR /app/morgan

# Run air for auto-reload with explicit config location
CMD ["air"]

# Stage 3: Production
FROM alpine:3.19 AS prod

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /go-morgan ./server
COPY --from=builder /app/morgan/config/config.example.json ./config/config.json

EXPOSE 8080

CMD ["./server", "serve"]

# Default to production
FROM prod
