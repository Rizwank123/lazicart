-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =========================
-- PRODUCTS
-- =========================
CREATE TABLE "public"."products" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),

    "name" varchar(255) NOT NULL,
    "description" text,

    "price" numeric(12,2) NOT NULL,
    "currency" varchar(10) NOT NULL DEFAULT 'INR',

    "category_id" uuid,
    "brand" varchar(255),

    "attributes" jsonb NOT NULL DEFAULT '{}'::jsonb,

    "is_active" bool NOT NULL DEFAULT true,

    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp,

    PRIMARY KEY ("id")
);

CREATE INDEX "idx_products_category_id"
ON "public"."products" ("category_id");

CREATE INDEX "idx_products_is_active"
ON "public"."products" ("is_active");

CREATE INDEX "idx_products_name_search"
ON "public"."products"
USING GIN (to_tsvector('english', "name"));


-- =========================
-- REVIEWS
-- =========================
CREATE TABLE "public"."reviews" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),

    "product_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,

    "rating" int NOT NULL CHECK ("rating" BETWEEN 1 AND 5),
    "comment" text,

    "is_approved" bool NOT NULL DEFAULT true,

    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp,

    PRIMARY KEY ("id"),

    CONSTRAINT "fk_reviews_product"
        FOREIGN KEY ("product_id")
        REFERENCES "public"."products" ("id")
        ON DELETE CASCADE
);

CREATE INDEX "idx_reviews_product_id"
ON "public"."reviews" ("product_id");

CREATE INDEX "idx_reviews_user_id"
ON "public"."reviews" ("user_id");


-- =========================
-- PRODUCT RATINGS (AGGREGATE)
-- =========================
CREATE TABLE "public"."product_ratings" (
    "product_id" uuid NOT NULL,

    "average_rating" numeric(3,2) NOT NULL DEFAULT 0,
    "total_reviews" int NOT NULL DEFAULT 0,

    PRIMARY KEY ("product_id"),

    CONSTRAINT "fk_ratings_product"
        FOREIGN KEY ("product_id")
        REFERENCES "public"."products" ("id")
        ON DELETE CASCADE
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS "public"."product_ratings";
DROP TABLE IF EXISTS "public"."reviews";
DROP TABLE IF EXISTS "public"."products";

-- +goose StatementEnd