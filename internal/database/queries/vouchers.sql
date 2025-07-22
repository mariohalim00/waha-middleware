-- name: GetOneVoucher :one
SELECT * FROM "vouchers"
WHERE "name" = $1;
