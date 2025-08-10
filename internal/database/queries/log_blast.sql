
-- name: CreateLogBlast :one
INSERT INTO "log_blast" (workflow_start, blast_start, blast_end, actual_blast, success_blast, failed_blast, raw_blast, non_existent_number) 
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateLogBlast :one
UPDATE "log_blast"
SET  blast_start = $2, blast_end = $3,
    actual_blast = $4, success_blast = $5, failed_blast = $6
WHERE id = $1
RETURNING *;