FROM golang:1.25-alpine3.22 AS builder

WORKDIR /autumn-2025/

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY ./cmd/ ./cmd/
COPY ./config/ ./config/
COPY ./internal/ ./internal/
COPY ./pkg/ ./pkg/

RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/server ./cmd/server


FROM alpine:3.22 AS service

ARG MIGRATE_VERSION=v4.19.0

WORKDIR /service/

RUN apk --no-cache add ca-certificates tzdata curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/$MIGRATE_VERSION/migrate.linux-amd64.tar.gz | tar xvz

COPY ./migrations/ ./migrations/
COPY --from=builder /autumn-2025/build/server ./

CMD ["sh", "-c", "./migrate -source file://./migrations/ -database postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable up && ./server"]
