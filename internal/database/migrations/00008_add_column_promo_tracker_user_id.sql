-- +goose Up
-- +goose StatementBegin
ALTER TABLE "promo_tracker"
ADD COLUMN "user_id" VARCHAR NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "promo_tracker"
DROP COLUMN "user_id";
-- +goose StatementEnd
