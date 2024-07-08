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

metrics-test-4: build
	SERVER_PORT=$(random unused-port)
	ADDRESS="localhost:${SERVER_PORT}"
	TEMP_FILE=$(random tempfile)
	metricstest -test.v -test.run=^TestIteration4$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-server-port=$SERVER_PORT \
		-source-path=.

metrics-test-5: build
	SERVER_PORT=$(random unused-port)
	ADDRESS="localhost:${SERVER_PORT}"
	TEMP_FILE=$(random tempfile)
	metricstest -test.v -test.run=^TestIteration5$ \
		-agent-binary-path=cmd/agent/agent \
		-binary-path=cmd/server/server \
		-server-port=$SERVER_PORT \
		-source-path=.


autotest-sprint-1: static-test metrics-test-1 metrics-test-2 metrics-test-3 metrics-test-4 metrics-test-5

test:
	go test -count=1 ./...


.PHONY: all