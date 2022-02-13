.PHONY: run
run:
	go run ./cmd/main 1 "localhost:8080"

.PHONY: run1
run1:
	go run ./cmd/main 1 "localhost:8081"

.PHONY: run2
run2:
	go run ./cmd/main 2 "localhost:8082"

.PHONY: run3
run3:
	go run ./cmd/main 3 "localhost:8083"