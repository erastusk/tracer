BIN_NAME=tracer

build:
	go build -o ./app/${BIN_NAME}

run: build
	./app/${BIN_NAME}

gen:
	go generate -v ./...

ex:
	./app/${BIN_NAME}

clean:
	go clean
	rm ./app/${BIN_NAME}

test:
	go test ./... -v -cover

docker:
	@echo "Define docker build/run commands here"
