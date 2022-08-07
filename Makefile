start:
	go run api/server.go

update:
	./scripts/fetch-data.sh

generate:
	cd api && go run github.com/99designs/gqlgen generate

build-db:
	docker build -t idea456/painmon-api-go/db ./db

start-db:
	docker run -it -p 6379:6379 idea456/painmon-api-go/db