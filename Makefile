VERSION ?= $(shell git describe --tags)
COMMIT  ?= $(shell git rev-parse HEAD)
LDFLAGS ?= -X main.version=${VERSION} -X main.commit=${COMMIT}

.PHONY: build
build:
	@echo "==> Building native binary"
	@GO111MODULE=on go build -mod=vendor -v -ldflags="${LDFLAGS}"

.PHONY: test
test:
	@echo "==> Running tests"
	@GO111MODULE=on go test -mod=vendor -v ./...

.PHONY: vendor
vendor:
	@echo "==> Updating vendored packages"
	@GO111MODULE=on go mod tidy
	@GO111MODULE=on go mod vendor

.PHONY: docker
docker:
	@echo "==> Building Docker image"
	docker build --build-arg=VERSION=${VERSION} --build-arg=COMMIT=${COMMIT} -t andrewsomething/github-download-stats .