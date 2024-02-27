setup:
	brew tap golangci/tap
	brew install golangci/tap/golangci-lint

lint:
	golangci-lint run

goimport:
	 goimports -d $(find . -type f -name '*.go' -not -path "./vendor/*")

compose_build:
	docker compose --project-name gofipe -f ./deployments/docker-compose.yaml up --build

test:
	set -a && . ./configs/dev.env
	${GOROOT}/bin/go test -v ./...

coverage:
	set -a && . ./configs/dev.env
	${GOROOT}/bin/go test -coverprofile=coverage.out ./...
	${GOROOT}/bin/go tool cover -html=coverage.out

build-mocks:
	${GOROOT}/bin/go generate ./...
