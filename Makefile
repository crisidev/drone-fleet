.PHONY: all deps build clean fmt vet test docker

EXECUTABLE ?= drone-fleet
IMAGE ?= crisidev/$(EXECUTABLE)
COMMIT ?= $(shell git rev-parse --short HEAD)

LDFLAGS = -X "main.buildCommit=$(COMMIT)"
PACKAGES = $(shell go list ./... | grep -v /vendor/)

all: deps build

deps: clean
	mkdir -p .build
	go get -u github.com/axw/gocov/gocov
	go get -u github.com/drone/drone-go/drone
	go get -u github.com/drone/drone-go/plugin
	go get -u github.com/op/go-logging
	wget -q http://repo.crisidev.org/amd64/fleetctl -O .build/fleetctl
	chmod +x .build/fleetctl

$(EXECUTABLE): $(wildcard *.go)
	go build -ldflags '-s -w $(LDFLAGS) --extldflags "-static"'

build: $(EXECUTABLE)
	mv $(EXECUTABLE) .build

fmt:
	go fmt $(PACKAGES)

vet:
	go vet $(PACKAGES)

test:
	@for PKG in $(PACKAGES); do go test -cover -coverprofile coverage.out $$PKG || exit 1; done;
	gocov convert coverage.out | gocov report

docker: all
	docker build --rm -t $(IMAGE) .

release: all vet test
	drone secure --repo ${IMAGE} --in secrets.yml

clean:
	rm -rf build
