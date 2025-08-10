-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "log_blast" (
	"id" UUID NOT NULL UNIQUE  DEFAULT uuid_generate_v4(),
	"workflow_start" TIMESTAMPTZ,
	"blast_start" TIMESTAMPTZ,
	"blast_end" TIMESTAMPTZ,
	"actual_blast" INTEGER,
	"success_blast" INTEGER,
	"failed_blast" INTEGER,
	"raw_blast" INTEGER,
	"non_existent_number" INTEGER,
	PRIMARY KEY("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "log_blast";
-- +goose StatementEnd
