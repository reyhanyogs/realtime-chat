postgresinit:
	docker run --name realtimechat-postgres -p 5433:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -d postgres:16-alpine

postgres:
	docker exec -it realtimechat-postgres psql --username=admin

createdb:
	docker exec -it realtimechat-postgres createdb --username=admin --owner=admin go-chat

dropdb:
	docker exec -it realtimechat-postgres dropdb go-chat

migrateup:
	migrate -path db/migrations -database "postgresql://admin:admin@localhost:5433/go-chat?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://admin:admin@localhost:5433/go-chat?sslmode=disable" -verbose down

mockuserservice:
	mockgen -package=mocks -destination domain/mocks/Service.go github.com/reyhanyogs/realtime-chat/domain Service

mockuserrepository:
	mockgen -package=mocks -destination domain/mocks/Repository.go github.com/reyhanyogs/realtime-chat/domain Repository

.PHONY: postgresinit postgres createdb dropdb migrateup migratedown mockuserservice mockuserrepository