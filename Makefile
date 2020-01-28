CGO    = 0
OUTPUT_FILE = dist/cfdd

VERSION ?= DEV

all: clean build

test:
	@go test ./...

clean:
	@go clean
	@rm $(OUTPUT_FILE) -f

build:
	@mkdir -p dist
	@CGO_ENABLED=$(CGO) go build -ldflags "-w -s -X main.version=${VERSION} -X main.date=$(shell date +'%y.%m.%dT%H:%M:%S')" \
	                             -o $(OUTPUT_FILE) \
								 main.go

release:
	@VERSION=$(VERSION) goreleaser --skip-publish --snapshot --rm-dist
