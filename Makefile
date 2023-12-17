compose_build:
	docker compose --project-name gofipe -f ./deployments/docker-compose.yaml up --build -d
