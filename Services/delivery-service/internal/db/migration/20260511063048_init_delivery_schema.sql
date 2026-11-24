-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS "public"."shipments" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "order_id" uuid UNIQUE NOT NULL,
    "user_id" uuid NOT NULL,
    "carrier" varchar(100),
    "tracking_number" varchar(160),
    "status" varchar(40) NOT NULL DEFAULT 'pending',
    "shipping_address" jsonb NOT NULL DEFAULT '{}'::jsonb,
    "estimated_delivery_at" timestamp,
    "delivered_at" timestamp,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);

CREATE INDEX IF NOT EXISTS "idx_shipments_user_id"
ON "public"."shipments" ("user_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."shipments";
-- +goose StatementEnd
