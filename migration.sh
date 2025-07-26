#!/bin/sh

migration_type=$1

if [ -z "$migration_type" ]; then
  echo "Usage: $0 <migration_type>"
  echo "Example: $0 up | down"
  exit 1
fi

goose -dir ./projects/internal/database/migrations postgres "host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable" $migration_type