protos:
	# Users
	protoc --go_out=./services/users/pkg \
	--go-grpc_out=require_unimplemented_servers=true:./services/users/pkg \
	./services/users/resources/users.proto

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
	GOOS=linux go build -tags musl -o ./services/users/build/users ./services/users/cmd/users/main.go

.PHONY: clean
clean:
	rm -rf build
	rm -rf ./services/users/build

.PHONY: build-docker
build-docker:
	docker build --no-cache -t users:local . --build-arg SERVICE=users

.PHONY: run-docker
run-docker:
	docker run -it --entrypoint sh users:local

.PHONY: generate-mocks
generate-mocks:
	go generate ./...