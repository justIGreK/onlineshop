postgres:
	docker run --name postgres -p 5555:5432 -e POSTGRESS_USER=root -e POSTGRES_PASSWORD=secret -d postgres
migrate up:
	migrate -path ./schema -database 'postgres://postgres:secret@localhost:5555/postgres?sslmode=disable' up  
linter: 
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1 run -v
migrate down:
	migrate -path ./schema -database 'postgres://postgres:secret@localhost:5555/postgres?sslmode=disable' down 
.PHONY: migrateup, migratedown