# Avito autumn 2025 test task

## Project start

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

I'm using `golangci-lint` as project linter.
CI pipeline placing at `.github/workflows/lint.yaml`.
You can also find configuration in `.golangci.yaml`. 

## Load testing

I'm using `siege` for load tests:

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
