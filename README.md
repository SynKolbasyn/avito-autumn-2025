# Avito autumn 2025 test task

- [Project start](#Project-start)
  - [Configuration](#Configuration)
  - [Services up](#Services-up)
- [Lint](#Lint)
- [Load testing](#Load-testing)

## Project start

### Configuration

If you want to configure service you must copy [`.env.example`](.env.example) to `.env` and change values for what you want.

```bash
cp .env.example .env
```

### Services up

You can use either

```bash
make up
```

or

```bash
docker compose up --build --pull always -d
```

this will up 2 containers: `postgres` and `service`

## Lint

I'm using [`golangci-lint`](https://github.com/golangci/golangci-lint) as project linter.
CI pipeline placing at [`.github/workflows/lint.yaml`](.github/workflows/lint.yaml).
You can also find configuration in [`.golangci.yaml`](.golangci.yaml). 

You can check project using either

```bash
make lint
```

or

```bash
golangci-lint run ./...
```

## Load testing

I'm using [`siege`](https://github.com/JoeDog/siege) for load tests:

```bash
siege -c 4 -b -t 10s http://localhost:8080/health
```

output:

```
Transactions:                  13043 hits
Availability:                 100.00 %
Elapsed time:                   9.65 secs
Data transferred:               0.63 MB
Response time:                  0.00 secs
Transaction rate:            1352.03 trans/sec
Throughput:                     0.07 MB/sec
Concurrency:                    3.15
Successful transactions:       13044
Failed transactions:               0
Longest transaction:            0.26
Shortest transaction:           0.00
```
