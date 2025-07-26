#!/bin/sh

migration_type=$1

if [ -z "$migration_type" ]; then
  echo "Usage: $0 <migration_type>"
  echo "Example: $0 up | down"
  exit 1
fi

if [ -z "$DB_HOST" ] || [ -z "$DB_USERNAME" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_DATABASE" ]; then
  echo "Error: One or more required environment variables (DB_HOST, DB_USERNAME, DB_PASSWORD, DB_DATABASE) are not set."
  exit 1
fi

goose -dir ./projects/internal/database/migrations postgres "host=${DB_HOST} user=${DB_USERNAME} password=${DB_PASSWORD} dbname=${DB_DATABASE} sslmode=disable" ${migration_type}