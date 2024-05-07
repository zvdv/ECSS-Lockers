#!/bin/sh
# This script is meant to be running inside the database container 
# specified in the docker-compose.yml file.

DB_ARGS=$(echo $DATABASE_URL | awk -F"[:/@?]" '{print "-h " $6 " -u " $4 " -p" $5}' | xargs)
DATABASE='locker-registration'

mysql -u $MYSQL_USER -p$MYSQL_PASSWORD $MYSQL_DATABASE < /home/schema.sql
