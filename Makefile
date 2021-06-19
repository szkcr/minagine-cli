.PHONY: run
run:
	go run .

.PHONY: build
build:
	go build -o minagine-cli .

.PHONY: dockerize
dockerize:
	docker build -t minagine-cli .
