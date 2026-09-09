FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -trimpath -ldflags="-s -w" -o /out/taskmanager ./cmd/api

FROM alpine:3.22

RUN addgroup -S app && adduser -S -G app app

COPY --from=builder /out/taskmanager /usr/local/bin/taskmanager

USER app
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/taskmanager"]