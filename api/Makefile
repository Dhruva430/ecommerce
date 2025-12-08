### CONFIGURATION ###
DB_URL = postgres://user:password@localhost:5432/mydatabase?sslmode=disable
MIGRATIONS_DIR = ./db/migrations


### COMMANDS ###
.PHONY: migrate-up migrate-down migrate-force reset create-migration

## Apply all up migrations
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose up

## Rollback last migration
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose down 

## Force set migration version (fix dirty state)
migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose force $(v)

## Reset DB: drop schema + apply migration again (⚠ deletes everything)
reset:
	docker exec -i ecommerce-postgres psql -U user -d mydatabase -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose up

## Create a new migration file: make create-migration name=add_discount
create-migration:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)
