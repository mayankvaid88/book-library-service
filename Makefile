unit_test:
	go test `go list ./... | grep -v './internal/integration_test'`
int_test:
	set -a && source .env && set +a && docker-compose up --detach db && sleep 5 && go test ./internal/integration_test && docker-compose down db
vet:
	go vet ./...
build_img:
	docker build . -t book-store-service
run:
	docker-compose up --detach
swag:
	swag init -g cmd/main.go
mocks:
	mockgen -source=internal/book/repository.go -destination=internal/mocks/repository_mock.go
	mockgen -source=internal/book/service.go -destination=internal/mocks/service_mock.go
