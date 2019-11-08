ifeq ($(OS),Windows_NT)
  SEPARATOR := &&
  SET := set
endif

build:
	$(eval DATE := $(shell go run scripts/getDate.go))
	$(eval GITHASH := $(shell git rev-parse --short HEAD))
	$(SET) CGO_ENABLED=1 $(SEPARATOR) go build -mod=vendor -v -ldflags "-X 'github.com/stashapp/stashdb/pkg/api.buildstamp=$(DATE)' -X 'github.com/stashapp/stashdb/pkg/api.githash=$(GITHASH)'"

install:
	packr2 install

clean:
	packr2 clean

# Regenerates GraphQL files and packr files
.PHONY: generate
generate:
	go generate

# Runs gofmt -w on the project's source code, modifying any files that do not match its style.
.PHONY: fmt
fmt:
	go fmt ./...

# Runs go vet on the project's source code.
.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	revive -config revive.toml -exclude ./vendor/...  ./...
