#!/bin/sh

case "$1" in
  "server")
    exec /app/server
    ;;
  "migrate-up")
  exec /app/migrate -path /app/migrations -database "$DB_CONNECTION" up
  ;;
  "seeder")
    exec /app/seeder "$2"
    ;;
  *)
    echo "Please choose between seeder,server or migrate-up"
    exit 1
    ;;
esac
