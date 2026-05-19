FROM golang:1.26-alpine AS builder

ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/ragforge ./apps/ragforge/cmd/

FROM alpine:3.22

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/ragforge /app/ragforge
COPY --from=builder /build/apps/ragforge/migrations /app/migrations
COPY --from=builder /build/apps/ragforge/config/config.yaml /app/config/config.yaml
COPY apps/ragforge/scripts/docker-entrypoint.sh /app/docker-entrypoint.sh

RUN chmod +x /app/docker-entrypoint.sh

WORKDIR /app
EXPOSE 8082

ENTRYPOINT ["/app/docker-entrypoint.sh"]
