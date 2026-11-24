-- +goose Up
-- +goose StatementBegin
-- =========================
-- PRODUCT IMAGES
-- =========================
CREATE TABLE "public"."product_images" (
    "id" uuid NOT NULL DEFAULT gen_random_uuid(),

    "product_id" uuid NOT NULL,
    "image_url" text NOT NULL,
    "thumbnail_url" text,

    "is_primary" bool NOT NULL DEFAULT false,
    "display_order" int NOT NULL DEFAULT 0,

    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp,

    PRIMARY KEY ("id"),

    CONSTRAINT "fk_product_images_product"
        FOREIGN KEY ("product_id")
        REFERENCES "public"."products" ("id")
        ON DELETE CASCADE
);

CREATE INDEX "idx_product_images_product_id"
ON "public"."product_images" ("product_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_product_images_product_id";
DROP TABLE IF EXISTS "public"."product_images";
-- +goose StatementEnd
