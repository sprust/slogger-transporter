dev-init:
	cp -i .env.example .env
	make docker-build

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-stop:
	docker-compose stop

restart:
	make stop
	make up

run:
	go run ./internal/server.go ${c}

build:
	go build -o ./bin/ internal/server.go
