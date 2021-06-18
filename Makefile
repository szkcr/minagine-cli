.PHONY: run
run:
	go run main.go

.PHONY: dockerize
dockerize:
	docker build -t automate-minagine .
