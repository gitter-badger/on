GOTOOLS = \
	github.com/golang/lint/golint

BIN    = $(GOPATH)/bin
GOLINT = $(BIN)/golint

VERSION=$(shell cat VERSION)
BUILD_TAGS?=autogen
BUILD=`git rev-parse HEAD`
GIT_COMMIT="$(shell git rev-parse --short HEAD)"
GIT_DIRTY="$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)"
GIT_DESCRIBE="$(shell git describe --tags --always)"
GIT_IMPORT="continuul.io/lsr/cmd/version"
LDFLAGS=-ldflags "-X ${GIT_IMPORT}.Version='${VERSION}' -X ${GIT_IMPORT}.GitCommit='${GIT_COMMIT}${GIT_DIRTY}' -X ${GIT_IMPORT}.GitDescribe='${GIT_DESCRIBE}'"

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.DEFAULT_GOAL := all

.PHONY: all
all:
	@go install --tags $(BUILD_TAGS) $(LDFLAGS) .

.PHONY: tools
tools:
	@go get -u -v $(GOTOOLS)

.PHONY: fmt
fmt:
	@gofmt -l -w $(SRC)

.PHONY: lint
lint:
	@go list ./... \
		| grep -v /vendor/ \
		| cut -d '/' -f 4- \
		| xargs -n1 \
			golint ;\

.PHONY: clean
clean:
	@go clean -i
