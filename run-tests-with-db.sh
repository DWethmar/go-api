#!/bin/bash

set -e
set -u

SCRIPTPATH="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

PQ_HOST="localhost"
PG_PORT=5433
PQ_DB_NAME="testdb"
PQ_USER="postgres"
PQ_PASSWORD="postgres"
DRIVER_NAME="postgres"     
CREATE_DB_SQL_FILE="${SCRIPTPATH}/scripts/sql/init.sql"     

# connect with: psql -h localhost -p 5433  -U postgres --password
docker run --name testing-postgres --rm -e POSTGRES_PASSWORD=$PQ_PASSWORD -d -p $PG_PORT:5432 postgres

echo "waiting for database to start..."

RETRIES=30
until PGPASSWORD=$PQ_PASSWORD psql -h localhost -p $PG_PORT -U $PQ_USER -c "select 1" > /dev/null 2>&1 || [ $RETRIES -eq 0 ]; do
  echo "Waiting for postgres server, $((RETRIES--)) remaining attempts..."
  sleep 1
done

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
PQ_DB_NAME=$PQ_DB_NAME \
PQ_PASSWORD=$PQ_PASSWORD \
DRIVER_NAME=$DRIVER_NAME \
CREATE_DB_SQL_FILE=$CREATE_DB_SQL_FILE \
go test -v ./...

