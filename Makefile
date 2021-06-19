.PHONY: run
run:
	go run .

.PHONY: build
build:
	go build -o automate-minagine .

.PHONY: dockerize
dockerize:
	docker build -t automate-minagine .
