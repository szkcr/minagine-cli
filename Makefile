.PHONY: run
run:
	go run ./src

.PHONY: build
build:
	go build -o minagine-cli ./src

.PHONY: dockerize
dockerize:
	docker build -t minagine-cli .
