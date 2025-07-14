-- +goose Up
ALTER TABLE "promo_tracker"
	ADD COLUMN "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = NOW();
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER promo_tracker_updated_at_trigger
BEFORE UPDATE ON "promo_tracker"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE UNIQUE INDEX "promo_tracker_hashed_string_unique"
	ON "promo_tracker" ("hashed_string");

CREATE INDEX "promo_tracker_index_0"
	ON "promo_tracker" ("hashed_string");

-- +goose Down
ALTER TABLE "promo_tracker"
	DROP COLUMN "updated_at";

DROP TRIGGER promo_tracker_updated_at_trigger ON "promo_tracker";

DROP FUNCTION update_updated_at_column();

DROP INDEX "promo_tracker_hashed_string_unique";

DROP INDEX "promo_tracker_index_0";
