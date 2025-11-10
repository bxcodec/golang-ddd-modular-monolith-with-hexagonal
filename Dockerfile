FROM golang:1.24.4-alpine AS builder

RUN apk add --no-cache make bash git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/engine .
COPY --from=builder /app/migrations ./migrations

RUN adduser -D -g '' appuser && chown -R appuser:appuser /app
USER appuser

EXPOSE 9090

ENTRYPOINT ["/app/engine"]
CMD ["rest", "--port", "9090"]

