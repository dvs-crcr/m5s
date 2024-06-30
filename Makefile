all: build

build:
	@echo "Building client/agent application..."
	go build -o cmd/server/server cmd/server/*.go
	go build -o cmd/agent/agent cmd/agent/*.go

run-server:
	go run ./cmd/server/.

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

autotest: static-test metrics-test-1

test:
	go test -count=1 ./...


.PHONY: all