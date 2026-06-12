DB_URL="postgresql://admin:secretpassword@localhost:5432/book_halal?sslmode=disable"

create_migration:
	migrate create -ext sql -dir migrations -seq $(name)

migrate_up:
	@echo "Waiting for database..."
	@sleep 5
	migrate -path migrations -database $(DB_URL) -verbose up
migrate_down:
	migrate -path migrations -database $(DB_URL) -verbose down 1
