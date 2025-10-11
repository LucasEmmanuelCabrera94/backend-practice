SHELL := /bin/zsh

.PHONY: run build clean

run:
	# Ejecuta con CGO deshabilitado para evitar problemas de enlace en macOS
	CGO_ENABLED=0 go run ./cmd

build:
	CGO_ENABLED=0 go build -o bin/server ./cmd

clean:
	rm -rf bin

docker-build:
	docker build -t auth-service .

docker-run:
	# Run detached and map port 8080
	docker run -d -p 8080:8080 --name auth-service auth-service || true

docker-stop:
	docker stop auth-service || true && docker rm auth-service || true
