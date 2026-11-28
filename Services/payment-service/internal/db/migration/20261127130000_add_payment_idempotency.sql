-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "public"."payment_idempotency_keys" (
    "idempotency_key" varchar(255) PRIMARY KEY,
    "request_hash" varchar(64) NOT NULL,
    "payment_id" uuid NOT NULL REFERENCES "public"."payments" ("id"),
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS "idx_payment_idempotency_payment_id"
ON "public"."payment_idempotency_keys" ("payment_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."payment_idempotency_keys";
-- +goose StatementEnd
