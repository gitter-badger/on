GOTOOLS = \
	github.com/golang/lint/golint \
	github.com/golang/dep/cmd/dep \
	github.com/mitchellh/gox

BIN    = $(GOPATH)/bin
GOLINT = $(BIN)/golint
CGO_ENABLED=0
VERSION=$(shell cat VERSION)
BUILD_TAGS?=autogen
BUILD=`git rev-parse HEAD`
GIT_COMMIT="$(shell git rev-parse --short HEAD)"
GIT_DIRTY="$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)"
GIT_DESCRIBE="$(shell git describe --tags --always)"
GIT_IMPORT="continuul.io/on/cmd/version"
LDFLAGS=-ldflags "-X ${GIT_IMPORT}.Version='${VERSION}' -X ${GIT_IMPORT}.GitCommit='${GIT_COMMIT}${GIT_DIRTY}' -X ${GIT_IMPORT}.GitDescribe='${GIT_DESCRIBE}'"

XC_ARCH ?= "amd64"
XC_OS ?= "darwin linux"
ifneq ($(strip $(CONTINUUL_DEV)),)
    XC_OS=$(shell go env GOOS)
    XC_ARCH=$(shell go env GOARCH)
endif

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
PACKAGES := go list ./... | grep -v /vendor | grep -v /out

.DEFAULT_GOAL := all

.PHONY: vendor
vendor:
	@dep ensure -v

.PHONY: all
all: vendor
	@gox -os=$(XC_OS) -arch=$(XC_ARCH) $(LDFLAGS) --tags $(BUILD_TAGS) -output "$(GOPATH)/bin/{{.OS}}_{{.Arch}}/on" .
	@go install --tags $(BUILD_TAGS) $(LDFLAGS) .

.PHONY: tools
tools:
	@go get -u -v $(GOTOOLS)

.PHONY: fmt
fmt:
	@gofmt -l -w $(SRC)

.PHONY: vet
vet:
	@go vet $(shell $(PACKAGES))

.PHONY: lint
lint:
	@golint $(shell $(PACKAGES))

.PHONY: docker
docker:
	@cp $(GOPATH)/bin/linux_amd64/on .
	@docker build -t continuul/on:$(VERSION) .
	@rm on

.PHONY: clean
clean:
	@go clean -i
	@rm -fr vendor

.PHONY: version
version:
	@echo $(VERSION)