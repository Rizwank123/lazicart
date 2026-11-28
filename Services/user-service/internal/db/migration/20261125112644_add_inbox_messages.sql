-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "public"."inbox_messages" (
    "event_id" uuid PRIMARY KEY,
    "event_type" varchar(150) NOT NULL,
    "status" varchar(30) NOT NULL DEFAULT 'processing',
    "attempts" int NOT NULL DEFAULT 1,
    "last_error" text,
    "processed_at" timestamp,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_inbox_messages_status"
ON "public"."inbox_messages" ("status", "created_at");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."inbox_messages";
-- +goose StatementEnd