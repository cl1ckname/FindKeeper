FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN make build

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bot .
CMD ["make", "run"]
