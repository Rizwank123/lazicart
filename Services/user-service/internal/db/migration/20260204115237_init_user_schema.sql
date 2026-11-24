-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "public"."users" (
    "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    "first_name" character varying(255) NOT NULL,
    "last_name" character varying(255) NOT NULL,
    "email" character varying(255) NOT NULL,
    "phone" character varying(255) NOT NULL,
    "is_active" boolean NOT NULL,
    "dob" timestamp without time zone NOT NULL,
    "address" jsonb NOT NULL DEFAULT '{}',
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	"updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."users";
-- +goose StatementEnd
