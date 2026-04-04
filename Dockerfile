FROM golang:1.25-alpine AS builder

WORKDIR /app
ENV GOPROXY=https://goproxy.io,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app ./cmd/find-keeper

FROM alpine:latest
WORKDIR /app
COPY --from=builder app .
CMD ["./app"]
