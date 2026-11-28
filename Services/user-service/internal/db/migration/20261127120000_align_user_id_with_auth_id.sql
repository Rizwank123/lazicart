-- +goose Up
-- +goose StatementBegin
-- Align the profile id with the auth-service id.
--
-- The auth-service puts the auth-service user id in the JWT subject claim,
-- while the commerce services (cart, order, payment, notification, delivery)
-- key their rows by the user id the client sends. Making users.id equal
-- auth_id gives the whole platform one canonical UUID per user.
--
-- No foreign key references users.id, so the primary key can be rewritten.
UPDATE users SET id = auth_id WHERE id <> auth_id AND auth_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- The previous id cannot be reconstructed, so the change is not reversible.
SELECT 1;
-- +goose StatementEnd
