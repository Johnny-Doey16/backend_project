# docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=toor -d postgres:16.0-alpine3.18
DB_URL=postgresql://root:toor@localhost:5434/diivix?sslmode=disable

air_init:
	air init

air_run:
	air

start_ps:
	docker start postgres16

start_redis:
	docker start redis

postgres:
	docker run --name postgres16 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=toor -d postgres:16.0-alpine3.18

redis:
	docker run --name redis -p 6379:6379 -d redis:7.0-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root authentication

dropdb:
	docker exec -it postgres16 dropdb authentication

migrate_init:
	migrate create -ext sql -dir db/migration -seq prayer_room

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

migratedown2:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 2

db_docs:
	dbdocs build docs/db.dbml

db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

sqlc_init:
	sqlc init

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

run:
	go run main.go

proto:
	rm -f pb/*.go
	rm -f docs/swagger/*.swagger.json
	protoc --proto_path=protos --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=church \
	protos/*.proto

evans:
	evans --host localhost --port 9901 -r repl

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test run environ air_init air_run start_redis redis start_ps db_docs db_schema proto evans