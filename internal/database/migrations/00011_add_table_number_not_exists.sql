-- +goose Up
-- +goose StatementBegin
CREATE TABLE "phone_number_not_exist" (
	"phone_number" VARCHAR(255) NOT NULL UNIQUE,
	"username" VARCHAR(255),
	"blast_id" UUID,
	PRIMARY KEY("phone_number")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "phone_number_not_exist";
-- +goose StatementEnd
