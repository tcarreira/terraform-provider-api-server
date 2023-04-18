default: testacc

REGISTRY=registry.terraform.io
NAMESPACE=tcarreira
NAME=apiserver
BINARY=terraform-provider-${NAME}
VERSION=$(shell git describe --tags --match='v*' 2> /dev/null | sed 's/^v//')
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

build:
	go build -o ${BINARY} .

install: build
	mkdir -p ~/.terraform.d/plugins/${REGISTRY}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	@rm -f ~/.terraform.d/plugins/${REGISTRY}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}
	mv ${BINARY} ~/.terraform.d/plugins/${REGISTRY}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
