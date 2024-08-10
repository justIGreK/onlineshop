migrate up:
	migrate -path ./schema -database 'postgres://postgres:secret@localhost:5555/postgres?sslmode=disable' up  
migrate down:
	migrate -path ./schema -database 'postgres://postgres:secret@localhost:5555/postgres?sslmode=disable' down 
.PHONY: migrateup, migratedown