# ──────────────────────────────────────────────────────────────
# Stage: build
# ──────────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS build

WORKDIR /src

# Cache dependencies
COPY services/api/go.mod services/api/go.sum ./
RUN go mod download

COPY services/api/ .

RUN CGO_ENABLED=0 GOOS=linux go build \
      -ldflags="-s -w" \
      -o /retrosnack-api \
      ./cmd/server

# ──────────────────────────────────────────────────────────────
# Stage: development (for docker-compose hot-reload)
# ──────────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS development

RUN go install github.com/air-verse/air@latest

WORKDIR /app
COPY services/api/go.mod services/api/go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]

# ──────────────────────────────────────────────────────────────
# Stage: production
# ──────────────────────────────────────────────────────────────
FROM alpine:3.21 AS production

RUN apk --no-cache add ca-certificates nginx tzdata

COPY infrastructure/nginx/nginx.conf /etc/nginx/nginx.conf
COPY --from=build /retrosnack-api /usr/local/bin/retrosnack-api

EXPOSE 80

CMD ["sh", "-c", "nginx && retrosnack-api"]
