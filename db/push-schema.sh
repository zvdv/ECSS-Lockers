#!/bin/sh

DB_ARGS=$(echo $DATABASE_URL | awk -F"[:/@?]" '{print "-h " $6 " -u " $4 " -p" $5}' | xargs)
DATABASE='locker-registration'

mysql $DB_ARGS "$@" $DATABASE  < db/schema.sql