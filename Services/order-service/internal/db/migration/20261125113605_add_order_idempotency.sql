-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "public"."order_idempotency_keys" (
    "idempotency_key" varchar(255) PRIMARY KEY,
    "request_hash" varchar(64) NOT NULL,
    "order_id" uuid NOT NULL REFERENCES "public"."orders" ("id"),
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_order_idempotency_order_id"
ON "public"."order_idempotency_keys" ("order_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."order_idempotency_keys";
-- +goose StatementEnd