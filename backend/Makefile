DB_URL := "postgresql://admin:secretpassword@localhost:5432/book_halal?sslmode=disable"

# ─── Detect OS ────────────────────────────────────────────────────────────────
# On Windows (cmd/PowerShell), OS variable is set to "Windows_NT".
# On Linux / WSL, it is empty or "Linux".
ifeq ($(OS),Windows_NT)
  SLEEP := powershell -Command Start-Sleep -Seconds
  RM    := powershell -Command Remove-Item -Recurse -Force
  MKDIR := powershell -Command New-Item -ItemType Directory -Force -Path
else
  SLEEP := sleep
  RM    := rm -rf
  MKDIR := mkdir -p
endif

# ─── Migrations ───────────────────────────────────────────────────────────────
.PHONY: create_migration migrate_up migrate_down

create_migration:
	migrate create -ext sql -dir migrations -seq $(name)

migrate_up:
	@echo "Waiting for database..."
	@$(SLEEP) 5
	migrate -path migrations -database $(DB_URL) -verbose up

migrate_down:
	migrate -path migrations -database $(DB_URL) -verbose down 1

# ─── Build ────────────────────────────────────────────────────────────────────
.PHONY: build run tidy

build:
	go build -o bin/app ./cmd/...

run:
	go run ./cmd/...

tidy:
	go mod tidy

# ─── Docker ───────────────────────────────────────────────────────────────────
.PHONY: up down

up:
	docker-compose up -d

down:
	docker-compose down
