-- name: GetAllTrackedPromos :many
SELECT * FROM "promo_tracker"
ORDER BY id DESC;

-- name: GetOneTrackedPromo :one
SELECT * FROM "promo_tracker"
WHERE "hashed_string" = $1;

-- name: CreateTrackedPromo :one
INSERT INTO "promo_tracker" (hashed_string, expired_at, created_at, claimed, user_name) 
VALUES($1, $2, now(), false, $3)
RETURNING *;

-- name: UpdateTrackedPromo :one
UPDATE "promo_tracker"
SET claimed = $2
WHERE "hashed_string" = $1
RETURNING *;