Comming soon!

Migration command
docker run -v /Users/jovid/Desktop/Dev/microCmp/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:test@localhost:5432/postgres?sslmode=disable" up