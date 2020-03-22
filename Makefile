.DEFAULT_GOAL := help

ENTRYPOINT_API = cmd/couchconnections-api
CONTAINER ?= couchconnections-api
PORT ?= 8925
GOLANGCI_LINT_VERSION = v1.21.0
PACKR_VERSION := 2.2.0

# Setup name variables for the package/tool
NAME := couchconnections
PKG := github.com/sebastianrosch/$(NAME)
BUILDINFOPKG := $(PKG)/pkg/build-info
ALL_PKGS := $(shell go list ./...)

# Set build info
BUILDUSER := $(shell whoami)
BUILDTIME := $(shell date -u '+%Y-%m-%d %H:%M:%S')

# Populate version variables
GITCOMMIT := $(shell git rev-parse --short HEAD)
VERSION := $(GITCOMMIT)
GITUNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(GITUNTRACKEDCHANGES),)
	GITCOMMIT := $(GITCOMMIT)-dirty
endif
GITBRANCH = $(shell git rev-parse --verify --abbrev-ref HEAD)
CTIMEVAR = -X '$(BUILDINFOPKG).Version=$(VERSION)' \
 		   -X '$(BUILDINFOPKG).Revision=$(GITCOMMIT)' \
		   -X '$(BUILDINFOPKG).Branch=$(GITBRANCH)' \
		   -X '$(BUILDINFOPKG).BuildUser=$(BUILDUSER)' \
		   -X '$(BUILDINFOPKG).BuildDate=$(BUILDTIME)'

# Note this is done differently in the Jenkinsfile as part of CI.
.PHONY: docker-build
docker-build: clean ## builds the docker image
	docker build --build-arg BUILD_VERSION=$(VERSION) -t $(CONTAINER):latest .

.PHONY: docker-push
docker-push: docker-build ## pushes the docker image
	docker push $(CONTAINER):latest

.PHONY: docker-run
docker-run: ## runs the docker image
	docker run --rm -i -p $(PORT):$(PORT) \
    $(CONTAINER):latest

.PHONY: run
run: ## runs the api
	@HOST=localhost go run -ldflags "$(CTIMEVAR)" ./$(ENTRYPOINT_API)

.PHONY: build
build: install-packr2 ## builds the api
	cd ./$(ENTRYPOINT_API) && ../../bin/packr2
	CGO_ENABLED=0 go build -ldflags "-w $(CTIMEVAR) -s -extldflags -static" ./$(ENTRYPOINT_API)
	cd ./$(ENTRYPOINT_API) && ../../bin/packr2 clean

.PHONY: lint
# Runs https://github.com/golangci/golangci-lint
# Check config at .golangci.yml
 lint: ## runs the go linter
	@echo "+ $@"
	@golangci-lint run ./...

.PHONY: swagger-check
swagger-check: ## validates the swagger spec for syntax
	@echo "+ $@"
	@swagger validate api/swagger/v1/service.swagger.json

.PHONY: test
test: ## runs the unit tests
	gotestsum --no-summary=output,skipped --format short-with-failures -- \
	-tags="unit" \
	-ldflags "-X github.com/sebastianrosch/couchconnections/pkg/build-info.Version=0.0.1 -X github.com/sebastianrosch/couchconnections/pkg/build-info.Branch=local -X github.com/sebastianrosch/couchconnections/pkg/build-info.Revision=local -X github.com/sebastianrosch/couchconnections/pkg/build-info.BuildDate=`date -u +%Y%m%d.%H%M%S` -X github.com/sebastianrosch/couchconnections/pkg/build-info.BuildUser=`whoami`" \
	-cover -coverprofile cover.out -covermode=count ./...

.PHONY: clean
clean:  ## cleans up any build binaries or packages
	find . -name '*-packr.go' -type f -exec rm "{}" \;
	find . -name 'packrd' -type d -empty -delete
	rm -rf bin

.PHONY: install-tools
install-tools: install-packr2 ## installs the required build tools
	brew bundle && \
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b /usr/local/bin/ v1.21.0 && \
	brew tap go-swagger/go-swagger && brew install go-swagger && \
	go install \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
		github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc \
		github.com/golang/protobuf/protoc-gen-go \
		gotest.tools/gotestsum && \
	GO111MODULE=on go get -u github.com/sixt/protoc-gen-jsonschema/cmd/protoc-gen-jsonschema && go install github.com/sixt/protoc-gen-jsonschema/cmd/protoc-gen-jsonschema && \
	GO111MODULE=on go get github.com/golang/mock/mockgen@v1.4.3 && \
	go get -d github.com/envoyproxy/protoc-gen-validate@v0.3.0 && cd $(GOPATH)/src/github.com/envoyproxy/protoc-gen-validate && GO111MODULE=off make build && cd -

.PHONY: install-packr2
install-packr2: bin/packr2 ## installs the packr tool
bin/packr2:
	echo "running install-packr2"
	eval $$(go tool dist env) && \
	  mkdir -p bin && cd bin && \
    curl -s -f -L "https://github.com/gobuffalo/packr/releases/download/v$(PACKR_VERSION)/packr_$(PACKR_VERSION)_$${GOOS}_$${GOARCH}.tar.gz" | tar xz packr2

.PHONY: generate
generate: ## generates the source code from the protobuf file
	cd rpc/couchconnections-api && protoc \
		-I=. \
		-I=/usr/local/include \
		-I=$(GOPATH)/src \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway \
		-I=$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  		-I=${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate \
		--go_out=plugins=grpc:. \
		--grpc-gateway_out=logtostderr=true:. \
		--swagger_out=logtostderr=true:../../api/swagger \
		--validate_out=lang=go:. \
		--doc_out=../../doc --doc_opt=markdown,README.md \
		--jsonschema_out=disallow_additional_properties:../../api/schema/v1 \
		v1/service.proto && \
	rm -f v1/service.pb.mc.go && \
	mockgen -source v1/service.pb.go -mock_names CouchConnectionsAPI=MockCouchConnectionsAPI -destination v1/service.pb.mc.go -package v1

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
