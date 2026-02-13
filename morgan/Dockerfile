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

# Stage 2: Run
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies if needed (e.g. ca-certificates, tzdata)
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /go-morgan ./server
# Copy config if needed, usually passed via env or mounted volume, but copying defaults is safe
COPY --from=builder /app/morgan/config/config.example.json ./config/config.json

EXPOSE 8080

CMD ["./server", "serve"]
