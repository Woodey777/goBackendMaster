run_postgres:
	docker run -d --rm --name pgdb -e POSTGRES_PASSWORD=admin -e POSTGRES_USER=admin -e PGDATA=/var/lib/postgresql/data/pgdata -p 60000:5432 -v .:/var/lib/postgresql/data postgres

stop_postgres:
	docker stop pgdb

clean_postgres:
	rm -rf pgdata

create_db:
	docker exec -it pgdb createdb --username=admin --owner=admin  bank_db

drop_db:
	docker exec -it pgdb dropdb -U admin bank_db

migrate_up:
	migrate -path db/migration -database "postgresql://admin:admin@localhost:60000/bank_db?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgresql://admin:admin@localhost:60000/bank_db?sslmode=disable" -verbose down

generate_sqlc:
	sqlc generate

test:
	go test -v -cover db/sqlc_files/*