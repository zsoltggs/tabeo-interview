.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test ./...

.PHONY: test-all
test-all:
	E2E_TEST_ENABLED=1 INTEGRATION_TEST_ENABLED=1 go test ./...

.PHONY: verify
verify: lint test-all

.PHONY: build
build:
	GOOS=linux go build -o ./services/bookings/build/bookings ./services/bookings/cmd/bookings/main.go

.PHONY: clean
clean:
	rm -rf build
	rm -rf ./services/bookings/build

.PHONY: build-docker
build-docker:
	docker build --no-cache -t bookings:local . --build-arg SERVICE=bookings

.PHONY: run-docker
run-docker:
	docker run -it --entrypoint sh bookings:local

.PHONY: generate-sqlc
generate-sqlc:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	sqlc generate -f ./services/bookings/resources/database/sqlc.yaml

.PHONY: generate-mocks
generate-mocks:
	go generate ./...

.PHONY: generate
generate: generate-sqlc generate-mocks
