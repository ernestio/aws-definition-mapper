install:
	go install -v

build:
	go build -v ./...

lint:
	golint ./...
	go vet ./...

test:
	go test -v ./... --cover
	gucumber

deps:
	go get github.com/nats-io/nats
	go get github.com/r3labs/binary-prefix
	go get github.com/r3labs/graph
	go get github.com/r3labs/workflow
	go get github.com/ernestio/crypto
	go get github.com/ernestio/ernest-config-client

dev-deps: deps
	go get github.com/golang/lint/golint
	go get github.com/smartystreets/goconvey/convey
	go get -u github.com/gucumber/gucumber/cmd/gucumber


clean:
	go clean

dist-clean:
	rm -rf pkg src bin

