# Common targets
all: build

build:
	@echo "Building client/agent application..."
	go build -o cmd/server/server cmd/server/*.go
	go build -o cmd/agent/agent cmd/agent/*.go

run-server:
	go run ./cmd/server/*.go

run-agent:
	go run ./cmd/server/*.go

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

autotest-sprint-1: static-test metrics-test-1 metrics-test-2 metrics-test-3 metrics-test-4 metrics-test-5

.PHONY: static-test metrics-test-1 metrics-test-2 metrics-test-3 metrics-test-4 metrics-test-5 autotest-sprint-1

# Sprint 2
metrics-test-6: build
	metricstest -test.v -test.run=^TestIteration6$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-server-port=4485 \
		-source-path=.

metrics-test-7: build
	metricstest -test.v -test.run=^TestIteration7$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-server-port=4485 \
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