-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "public"."outbox_events" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "event_id" uuid NOT NULL UNIQUE,
    "aggregate_type" varchar(100) NOT NULL,
    "aggregate_id" uuid NOT NULL,
    "event_type" varchar(150) NOT NULL,
    "event_version" int NOT NULL,
    "payload" jsonb NOT NULL,
    "status" varchar(30) NOT NULL DEFAULT 'pending',
    "attempts" int NOT NULL DEFAULT 0,
    "available_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "published_at" timestamp,
    "last_error" text,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_order_outbox_pending"
ON "public"."outbox_events" ("status", "available_at", "created_at");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."outbox_events";
-- +goose StatementEnd
