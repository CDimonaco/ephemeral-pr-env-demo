#!/bin/sh

case "$1" in
  "server")
    exec /app/server
    ;;
  "seeder")
    exec /app/seeder "$2"
    ;;
  *)
    echo "Please choose between seeder or server"
    exit 1
    ;;
esac
