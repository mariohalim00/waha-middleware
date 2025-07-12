-- +goose Up
CREATE TABLE "promo_tracker" (
	"id" BIGSERIAL NOT NULL UNIQUE,
	"hashed_string" VARCHAR(255) NOT NULL,
	"expired_at" TIMESTAMPTZ NOT NULL,
	"created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	"claimed" BOOLEAN NOT NULL DEFAULT false,
	"user_name" VARCHAR(255) NOT NULL,
	PRIMARY KEY("id")
);

-- +goose Down
DROP TABLE "promo_tracker";