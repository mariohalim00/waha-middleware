
-- name: CreateLogBlast :one
INSERT INTO "log_blast" (workflow_start, blast_start, blast_end, actual_blast, success_blast, failed_blast, raw_blast, non_existent_number) 
VALUES($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateLogBlast :one
UPDATE "log_blast"
SET workflow_start = $2, blast_start = $3, blast_end = $4,
    actual_blast = $5, success_blast = $6, failed_blast = $7, raw_blast = $8, non_existent_number = $9
WHERE id = $1
RETURNING *;