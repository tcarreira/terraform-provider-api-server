default: testacc

REGISTRY=registry.terraform.io
NAMESPACE=tcarreira
NAME=apiserver
BINARY=terraform-provider-${NAME}
VERSION=$(shell git describe --tags --match='v*' 2> /dev/null | sed 's/^v//')
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)

test-server.pid:
	@echo "Starting test api-server"
	@API_PORT=18080 api-server . >> .testserver.log 2>&1 & echo $$! > .testserver.pid

stop-test-server.pid:
ifneq ("$(wildcard .testserver.pid)","")
	kill -9 $(shell cat .testserver.pid) ; rm .testserver.pid
endif

run-testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Run acceptance tests
.PHONY: testacc
testacc: stop-test-server.pid test-server.pid run-testacc stop-test-server.pid
	@echo "==> Tests finished"

build:
	go build -o ${BINARY} .

install: build
	mkdir -p ~/.terraform.d/plugins/${REGISTRY}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	@rm -f ~/.terraform.d/plugins/${REGISTRY}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}
	mv ${BINARY} ~/.terraform.d/plugins/${REGISTRY}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
