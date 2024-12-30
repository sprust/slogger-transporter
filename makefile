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
	go run main.go ${c}

run-serve-grpc:
	go run main.go serve:grpc

build:
	go build -o ./bin/ main.go

test:
	go test ./...

test-detail:
	go test -v ./...

grpc-generate:
	protoc --go_out=./internal/api/grpc/gen/ \
		--go-grpc_out=./internal/api/grpc/gen/ \
			./internal/api/grpc/proto/*.proto

bin-serve-grpc:
	./bin/main serve:grpc
