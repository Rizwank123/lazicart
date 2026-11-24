-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS "public"."inventory_items" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "product_id" uuid UNIQUE NOT NULL,
    "sku" varchar(120) UNIQUE,
    "quantity_available" int NOT NULL DEFAULT 0 CHECK ("quantity_available" >= 0),
    "quantity_reserved" int NOT NULL DEFAULT 0 CHECK ("quantity_reserved" >= 0),
    "reorder_level" int NOT NULL DEFAULT 0 CHECK ("reorder_level" >= 0),
    "is_active" boolean NOT NULL DEFAULT true,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);

CREATE INDEX IF NOT EXISTS "idx_inventory_items_product_id"
ON "public"."inventory_items" ("product_id");

CREATE TABLE IF NOT EXISTS "public"."inventory_reservations" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "order_id" uuid NOT NULL,
    "product_id" uuid NOT NULL,
    "quantity" int NOT NULL CHECK ("quantity" > 0),
    "status" varchar(30) NOT NULL DEFAULT 'reserved',
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp,
    UNIQUE ("order_id", "product_id")
);

CREATE INDEX IF NOT EXISTS "idx_inventory_reservations_order_id"
ON "public"."inventory_reservations" ("order_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."inventory_reservations";
DROP TABLE IF EXISTS "public"."inventory_items";
-- +goose StatementEnd
