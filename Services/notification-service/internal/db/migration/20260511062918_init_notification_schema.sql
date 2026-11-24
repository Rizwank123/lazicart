-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS "public"."notifications" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id" uuid,
    "channel" varchar(30) NOT NULL DEFAULT 'email',
    "type" varchar(60) NOT NULL,
    "recipient" varchar(255),
    "subject" varchar(255),
    "body" text NOT NULL,
    "status" varchar(30) NOT NULL DEFAULT 'pending',
    "metadata" jsonb NOT NULL DEFAULT '{}'::jsonb,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);

CREATE INDEX IF NOT EXISTS "idx_notifications_user_id"
ON "public"."notifications" ("user_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."notifications";
-- +goose StatementEnd
