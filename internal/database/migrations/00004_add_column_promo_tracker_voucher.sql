-- +goose Up
-- +goose StatementBegin
ALTER TABLE "promo_tracker"
    ADD COLUMN "voucher" VARCHAR NULL,
    ADD COLUMN "claimed_at" TIMESTAMPTZ DEFAULT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "promo_tracker"
    DROP COLUMN "voucher",
    DROP COLUMN "claimed_at";
-- +goose StatementEnd
