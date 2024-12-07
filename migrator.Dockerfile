FROM golang:1.23.3 AS builder

WORKDIR /db-migrations

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrator

FROM scratch

COPY --from=builder /db-migrations/schema/migrations /migrations
COPY --from=builder /db-migrations/migrate /migrate

CMD ["/migrate"]
