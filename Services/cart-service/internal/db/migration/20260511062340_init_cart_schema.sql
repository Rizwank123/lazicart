-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS "public"."carts" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id" uuid UNIQUE NOT NULL,
    "status" varchar(30) NOT NULL DEFAULT 'active',
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);

CREATE TABLE IF NOT EXISTS "public"."cart_items" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "cart_id" uuid NOT NULL REFERENCES "public"."carts" ("id") ON DELETE CASCADE,
    "product_id" uuid NOT NULL,
    "quantity" int NOT NULL CHECK ("quantity" > 0),
    "unit_price" numeric(12,2) NOT NULL DEFAULT 0,
    "currency" varchar(10) NOT NULL DEFAULT 'INR',
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp,
    UNIQUE ("cart_id", "product_id")
);

CREATE INDEX IF NOT EXISTS "idx_cart_items_cart_id"
ON "public"."cart_items" ("cart_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."cart_items";
DROP TABLE IF EXISTS "public"."carts";
-- +goose StatementEnd
