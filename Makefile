APP=iam-server

run:
	go run ./cmd/iam-server

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

race:
	go test -race ./...

tidy:
	go mod tidy

build:
	mkdir -p bin
	go build -o bin/$(APP) ./cmd/iam-server

clean:
	rm -rf bin