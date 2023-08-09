postgres:
	docker run --rm -d --name postgres15 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=1079 -p 5432:5432 postgres:15

createdb:
	docker exec -it postgres15 createdb --username=postgres --owner=postgres ecommerce

migrateup:
	migrate -path pkg/db/migration -database "postgresql://postgres:1079@localhost:5432/ecommerce?sslmode=disable" -verbose up

migratedown:
	migrate -path  pkg/db/migration -database "postgresql://postgres:1079@localhost:5432/ecommerce?sslmode=disable" -verbose down

dropdb:
	docker exec -it postgres15 dropdb -U postgres ecommerce

.PHONY: postgres createdb dropdb