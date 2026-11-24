-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS "public"."payments" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "order_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "amount" numeric(12,2) NOT NULL,
    "currency" varchar(10) NOT NULL DEFAULT 'INR',
    "provider" varchar(60) NOT NULL,
    "status" varchar(30) NOT NULL DEFAULT 'pending',
    "transaction_id" varchar(255),
    "failure_reason" text,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);

CREATE INDEX IF NOT EXISTS "idx_payments_order_id"
ON "public"."payments" ("order_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."payments";
-- +goose StatementEnd
