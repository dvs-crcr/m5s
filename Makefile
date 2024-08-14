# Common targets
all: build

build:
	@echo "Building client/agent application..."
	go build -o cmd/server/server cmd/server/*.go
	go build -o cmd/agent/agent cmd/agent/*.go

run-server:
	go run ./cmd/server/*.go -a localhost:44985 -i 10 -f tmp/storage

run-agent:
	go run ./cmd/agent/*.go -a localhost:44985

.PHONY: all build run-server run-agent

# Testing targets
test:
	go test -v -count=1 ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

lint:
	golangci-lint run ./...

.PHONY: test coverage lint

# Autotests
static-test:
	go vet -vettool=$(which statictest) ./...

# Sprint 1
metrics-test-1: build
	metricstest -test.v -test.run=^TestIteration1$ \
                -binary-path=cmd/server/server

metrics-test-2: build
	metricstest -test.v -test.run=^TestIteration2[AB]*$ \
                -source-path=. \
                -agent-binary-path=cmd/agent/agent

metrics-test-3: build
	metricstest -test.v -test.run=^TestIteration3[AB]*$ \
                -source-path=. \
                -agent-binary-path=cmd/agent/agent \
                -binary-path=cmd/server/server

metrics-test-4: build
	metricstest -test.v -test.run=^TestIteration4$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-server-port=4485 \
		-source-path=.

metrics-test-5: build
	metricstest -test.v -test.run=^TestIteration5$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-server-port=4485 \
		-source-path=.

autotest-sprint-1: static-test metrics-test-2 metrics-test-3 metrics-test-4 metrics-test-5

.PHONY: static-test metrics-test-1 metrics-test-2 metrics-test-3 metrics-test-4 metrics-test-5 autotest-sprint-1

# Sprint 2
metrics-test-6: build
	metricstest -test.v -test.run=^TestIteration6$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-server-port=4485 \
		-source-path=.

metrics-test-7: build
	export ADDRESS=localhost:8080; \
	metricstest -test.v -test.run=^TestIteration7$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-server-port=8080 \
		-source-path=.

metrics-test-8: build
	metricstest -test.v -test.run=^TestIteration8$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-server-port=4485 \
		-source-path=.

metrics-test-9: build
	metricstest -test.v -test.run=^TestIteration9$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-file-storage-path=tmp/tmp-storage \
		-server-port=4485 \
		-source-path=.

autotest-sprint-2: static-test metrics-test-6 metrics-test-7 metrics-test-8 metrics-test-9

.PHONY: metrics-test-6 metrics-test-7 metrics-test-8 metrics-test-9

# Sprint 3
metrics-test-10: build
	export ADDRESS=localhost:4485; \
	metricstest -test.v -test.run=^TestIteration10[AB]$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' \
		-server-port=4485 \
		-source-path=.

metrics-test-11: build
	export ADDRESS=localhost:4485; \
	metricstest -test.v -test.run=^TestIteration11$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' \
		-server-port=4485 \
		-source-path=.

metrics-test-12: build
	export ADDRESS=localhost:4485; \
	metricstest -test.v -test.run=^TestIteration12$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' \
		-server-port=4485 \
		-source-path=.

metrics-test-13: build
	export ADDRESS=localhost:4485; \
	metricstest -test.v -test.run=^TestIteration13$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' \
		-server-port=4485 \
		-source-path=.

autotest-sprint-3: static-test metrics-test-10 metrics-test-11 metrics-test-12 metrics-test-13

.PHONY: metrics-test-10 metrics-test-11 metrics-test-12 metrics-test-13 autotest-sprint-3

# Sprint 4
metrics-test-14: build
	export ADDRESS=localhost:4485; \
	metricstest -test.v -test.run=^TestIteration14$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-database-dsn='postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable' \
		-key="${TEMP_FILE}" \
		-server-port=4485 \
		-source-path=.

autotest-sprint-4: static-test metrics-test-14

.PHONY: metrics-test-14 autotest-sprint-4

