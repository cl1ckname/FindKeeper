FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o app ./cmd/find-keeper

FROM alpine:latest
WORKDIR /app
COPY --from=builder app .
CMD ["./app"]
