compose_build:
	docker compose --project-name gofipe -f ./deployments/docker-compose.yaml up --build -d

test-commit:
	pre-commit run --all-files

setup:
	brew install pre-commit yamllint
