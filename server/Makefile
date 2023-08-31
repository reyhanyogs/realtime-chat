postgresinit:
	docker run --name postgres16 -p 5433:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=admin -d postgres

postgres:
	docker exec -it postgres16 psql --username=postgres

createdb:
	docker exec -it postgres16 createdb --username=postgres --owner=postgres go-chat

dropdb:
	docker exec -it postgres16 dropdb go-chat

migrateup:
	migrate -path db/migrations -database "postgresql://postgres:admin@localhost:5433/go-chat?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:admin@localhost:5433/go-chat?sslmode=disable" -verbose down

.PHONY: postgresinit postgres createdb dropdb migrateup migratedown