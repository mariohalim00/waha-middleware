-- +goose Up
-- +goose StatementBegin
ALTER TABLE "vouchers"
ADD COLUMN "topup_deadline" INTEGER NULL,
ADD COLUMN "promo_text_template" TEXT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "vouchers"
DROP COLUMN "topup_deadline",
DROP COLUMN "promo_text_template";
-- +goose StatementEnd
