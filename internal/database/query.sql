-- name: GetAllTrackedPromos :many
SELECT * FROM "promo_tracker"
ORDER BY id DESC;
