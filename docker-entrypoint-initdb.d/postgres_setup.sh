#!/bin/bash

cd /migrations

for i in *.up.sql; do
    [ -f "$i" ] || break
    echo "Running $i";
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f $i
done

