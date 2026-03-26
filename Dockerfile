FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o battery-agent .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/battery-agent .

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s \
  CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["./battery-agent"]