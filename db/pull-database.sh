#!/bin/sh

DB_ARGS=$(echo $DATABASE_URL | awk -F"[:/@?]" '{print "-h " $6 " -u " $4 " -p" $5}' | xargs)
DATABASE='locker-registration'

# strip compatibility comments, we control the version anyways
mysqldump $DB_ARGS "$@" $DATABASE | grep -v '^/\*!.*\*/;'