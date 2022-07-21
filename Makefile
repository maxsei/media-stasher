all: build
build:
	echo "Building binary"
	go build -o ./bin/out *.go
	npm run build
run: build
	echo "Running binary"
	./bin/out
fmt:
	go run golang.org/x/tools/cmd/goimports -w .
	go fmt -w .
