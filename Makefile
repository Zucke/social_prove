default: build

build:
	go build -o app ./cmd/draid

serve: build
	./app && rm app

test:
	go test -v ./...

start:
	@if [ ! -f .env  ]; then\
		cp .env.example .env;\
	fi
	docker-compose -f docker-compose-local.yml up -d;

stop:
	docker-compose -f docker-compose-local.yml down
