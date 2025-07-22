-- +goose Up
-- +goose StatementBegin
ALTER TABLE "vouchers" ADD COLUMN "promo_duration_hours" INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "vouchers" DROP COLUMN "promo_duration_hours";
-- +goose StatementEnd
