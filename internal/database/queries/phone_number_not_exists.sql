-- name: CreatePhoneNumberNotExist :one
INSERT INTO "phone_number_not_exist" ("phone_number", "username", "blast_id")
VALUES ($1, $2, $3)
RETURNING *;