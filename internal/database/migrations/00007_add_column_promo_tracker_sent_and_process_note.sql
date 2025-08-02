-- +goose Up
-- +goose StatementBegin
ALTER TABLE "promo_tracker"
ADD COLUMN "sent_to_tm" BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN "process_note" TEXT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "promo_tracker"
DROP COLUMN "sent_to_tm",
DROP COLUMN "process_note";
-- +goose StatementEnd
