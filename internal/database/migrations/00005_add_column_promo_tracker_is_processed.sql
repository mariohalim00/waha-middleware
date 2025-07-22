-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "promo_tracker"
ADD COLUMN "is_processed" BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "promo_tracker"
DROP COLUMN "is_processed";
-- +goose StatementEnd
