migrate up:
	migrate -path ./schema -database 'postgres://postgres:secret@localhost:5555/shopdb?sslmode=disable' up  
migrate down:
	migrate -path ./schema -database 'postgres://postgres:secret@localhost:5555/shopdb?sslmode=disable' down 
.PHONY: migrateup, migratedown