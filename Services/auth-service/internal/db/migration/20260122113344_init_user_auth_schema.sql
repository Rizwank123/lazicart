-- +goose Up
-- +goose StatementBegin
CREATE TABLE "public"."user_auth" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	"email" VARCHAR(255) UNIQUE NOT NULL,
	"phone" VARCHAR(255) UNIQUE NOT NULL,
	"password" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	"updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL,
    "is_active" BOOLEAN DEFAULT TRUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."user_auth";
-- +goose StatementEnd
