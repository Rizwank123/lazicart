-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS "public"."orders" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id" uuid NOT NULL,
    "status" varchar(30) NOT NULL DEFAULT 'pending',
    "total" numeric(12,2) NOT NULL DEFAULT 0,
    "currency" varchar(10) NOT NULL DEFAULT 'INR',
    "shipping_address" jsonb NOT NULL DEFAULT '{}'::jsonb,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);

CREATE INDEX IF NOT EXISTS "idx_orders_user_id"
ON "public"."orders" ("user_id");

CREATE TABLE IF NOT EXISTS "public"."order_items" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "order_id" uuid NOT NULL REFERENCES "public"."orders" ("id") ON DELETE CASCADE,
    "product_id" uuid NOT NULL,
    "quantity" int NOT NULL CHECK ("quantity" > 0),
    "unit_price" numeric(12,2) NOT NULL,
    "currency" varchar(10) NOT NULL DEFAULT 'INR',
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."order_items";
DROP TABLE IF EXISTS "public"."orders";
-- +goose StatementEnd
