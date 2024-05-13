test:
	go test -v -count=1 ./...

migrate-up:
	migrate -path api/migrations -database "postgres://postgres:tdepassword@localhost:5432/auth_api?sslmode=disable" up

migrate-down:
	migrate -path api/migrations -database "postgres://postgres:tdepassword@localhost:5432/auth_api?sslmode=disable" down
