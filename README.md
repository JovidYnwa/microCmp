Comming soon!


for migration we should create folder (usualy is migrations)
then we should create file with the follog format id_action_table_name.up.sql (1481574547_create_users_table.up.sql)


Migration command
docker run -v /Users/jovid/Desktop/Dev/microCmp/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:test@localhost:5432/postgres?sslmode=disable" up

--Windows
docker run -v C:/Users/dzhovid.nurov/Desktop/dev/MicroCmp/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:test@localhost:5432/postgres?sslmode=disable" up

docker run -v C:/Users/dzhovid.nurov/Desktop/dev/MicroCmp/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "postgresql://postgres:test@localhost:5432/postgres?sslmode=disable" goto 3


docker-compose exec web go run scripts/seed.go
