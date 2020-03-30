#!/bin/bash

set -e
set -u

PQ_HOST="localhost"
PG_PORT=5433
PQ_USER="postgres"
PQ_PASSWORD="postgres"
DRIVER_NAME="postgres"     

SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

POSTGRES_MULTIPLE_DATABASES=test_one,test_two

# connect with: psql -h localhost -p 5433  -U postgres --password
docker run --name testing-postgres --rm -e POSTGRES_PASSWORD=$PQ_PASSWORD -d -p $PG_PORT:5432 postgres

echo "waiting for database to start..."

RETRIES=30
until PGPASSWORD=$PQ_PASSWORD psql -h localhost -p $PG_PORT -U $PQ_USER -c "select 1" > /dev/null 2>&1 || [ $RETRIES -eq 0 ]; do
  echo "Waiting for postgres server, $((RETRIES--)) remaining attempts..."
  sleep 1
done

function create_user_and_database() {
	local database=$1
	echo "  Creating user and database '$database' on $PG_PORT for user '$database' "
	PGPASSWORD=$PQ_PASSWORD psql -h localhost -v ON_ERROR_STOP=1 -p $PG_PORT -U $PQ_USER <<-EOSQL
	    CREATE USER $database WITH PASSWORD '$database';
	    CREATE DATABASE $database;
	    GRANT ALL PRIVILEGES ON DATABASE $database TO $database;
		GRANT ALL PRIVILEGES ON DATABASE $database TO postgres;

EOSQL
    PGPASSWORD=$PQ_PASSWORD psql -h localhost -v ON_ERROR_STOP=1 -d $database -p $PG_PORT -U $PQ_USER -f ${SCRIPTPATH}/scripts/sql/init.sql
}

if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
	echo "Multiple database creation requested: $POSTGRES_MULTIPLE_DATABASES"
	for db in $(echo $POSTGRES_MULTIPLE_DATABASES | tr ',' ' '); do
		create_user_and_database $db
	done
	echo "Multiple databases created"
fi

clean_up () {
    ARG=$?
    echo "> clean_up with exitcode: $ARG"
    docker stop testing-postgres
    exit $ARG
}
trap clean_up EXIT

PQ_HOST=$PQ_HOST \
PG_PORT=$PG_PORT  \
PQ_USER=$PQ_USER \
PQ_PASSWORD=$PQ_PASSWORD \
DRIVER_NAME=$DRIVER_NAME \
go test -run 'Intergration' -v ./...

