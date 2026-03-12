FROM golang:1.26-alpine AS build

WORKDIR /src

COPY services/api/go.mod services/api/go.sum ./
RUN go mod download

COPY services/api/ .

RUN CGO_ENABLED=0 GOOS=linux go build \
      -ldflags="-s -w" \
      -o /retrosnack-api \
      ./cmd/server

FROM golang:1.26-alpine AS development

RUN go install github.com/air-verse/air@latest

WORKDIR /app
COPY services/api/go.mod services/api/go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]

FROM alpine:3.21 AS production

RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -S app && adduser -S app -G app

COPY --from=build /retrosnack-api /usr/local/bin/retrosnack-api

USER app

EXPOSE 8080

CMD ["retrosnack-api"]
