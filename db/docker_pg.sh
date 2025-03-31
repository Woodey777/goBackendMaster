docker run -d --rm \
	--name pgdb \
	-e POSTGRES_PASSWORD=admin \
	-e POSTGRES_USER=admin \
	-e POSTGRES_DB=bank_db \
	-e PGDATA=/var/lib/postgresql/data/pgdata \
    -p 60000:5432 \
	-v .:/var/lib/postgresql/data \
	postgres
