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

test:
	go test ./...

test-detail:
	go test -v ./...

grpc-generate:
	protoc --go_out=./internal/grpc/gen/ \
		--go-grpc_out=./internal/grpc/gen/ \
			./internal/grpc/proto/ping_pong.proto

bin-server:
	./bin/server
