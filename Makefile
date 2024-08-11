postgres:
	docker run --name postgres -p 5555:5432 -e POSTGRESS_USER=root -e POSTGRES_PASSWORD=secret -d postgres
migrate up:
	migrate -path ./schema -database 'postgres://postgres:secret@localhost:5555/postgres?sslmode=disable' up  
migrate down:
	migrate -path ./schema -database 'postgres://postgres:secret@localhost:5555/postgres?sslmode=disable' down 
.PHONY: migrateup, migratedown