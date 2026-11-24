-- +goose Up
-- +goose StatementBegin
ALTER TABLE "public"."users" ADD COLUMN "auth_id" UUID NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "public"."users" DROP COLUMN "auth_id";
-- +goose StatementEnd
