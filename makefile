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
	go run ./cmd/transporter/main.go ${c}

run-start:
	go run ./cmd/transporter/main.go start

run-stop:
	go run ./cmd/transporter/main.go manage stop

run-serve-grpc:
	go run ./cmd/transporter/main.go serve:grpc

run-queue-listen:
	go run ./cmd/transporter/main.go queue:listen

build:
	CGO_ENABLED=0 GOOS=linux go build -v -a -o ./bin/ ./cmd/transporter/main.go \
		&& chmod +x ./bin/strans

test:
	go test ./...

test-detail:
	go test -v ./...

grpc-generate:
	protoc --go_out=./internal/api/grpc/gen/ \
		--go-grpc_out=./internal/api/grpc/gen/ \
			./internal/api/grpc/proto/*.proto

bin-start:
	./bin/strans start

bin-serve-grpc:
	./bin/strans serve:grpc

bin-queue-listen:
	./bin/strans queue:listen
