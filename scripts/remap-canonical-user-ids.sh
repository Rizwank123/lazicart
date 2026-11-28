#!/usr/bin/env bash
# Remap the user_id columns in the commerce databases to the canonical user id.
#
# Run this BEFORE the user-service migration
# 20261127120000_align_user_id_with_auth_id.sql, while the users table still
# has both its original id and the auth_id.
#
# The user-service owns the mapping (old users.id -> auth_id). Every other
# service keys its rows by that same user id, so each commerce row has to be
# updated to the canonical id before the primary key is rewritten.
set -euo pipefail

DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-root}"
DB_HOST="${DB_HOST:-localhost}"
export PGPASSWORD="$DB_PASSWORD"

psql_q() {
  psql -h "$DB_HOST" -U "$DB_USER" -d "$1" -tAc "$2"
}

echo "Reading the id mapping from the user-service..."
MAPPING=$(psql_q "user-service" "SELECT id::text || ' ' || auth_id::text
  FROM users
  WHERE auth_id IS NOT NULL AND id <> auth_id;" || true)

if [ -z "$MAPPING" ]; then
  echo "No user ids need remapping."
  exit 0
fi

echo "Mapping:"
echo "$MAPPING" | sed 's/^/  /'

# The commerce services are keyed by user_id in these tables.
TABLES=("cart:carts" "orders:orders" "payment:payments" "notification:notifications" "delivery:shipments")

while read -r old_id new_id; do
  [ -z "${old_id:-}" ] && continue
  echo "  remapping $old_id -> $new_id"

  # carts.user_id is UNIQUE, so a cart already stored under the canonical id
  # has to be merged into the source cart instead of colliding with it.
  source_cart=$(psql_q "cart" "SELECT id::text FROM carts WHERE user_id = '$old_id';" || true)
  target_cart=$(psql_q "cart" "SELECT id::text FROM carts WHERE user_id = '$new_id';" || true)
  if [ -n "$source_cart" ] && [ -n "$target_cart" ]; then
    echo "    merging cart $target_cart into $source_cart"
    psql_q "cart" "UPDATE cart_items SET cart_id = '$source_cart' WHERE cart_id = '$target_cart';" >/dev/null
    psql_q "cart" "DELETE FROM carts WHERE id = '$target_cart';" >/dev/null
  fi

  for spec in "${TABLES[@]}"; do
    db="${spec%%:*}"
    table="${spec##*:}"
    psql_q "$db" "UPDATE $table SET user_id = '$new_id' WHERE user_id = '$old_id';" >/dev/null
  done
done <<< "$MAPPING"

echo "Done. Apply the user-service migration next:"
echo "  make migrate SERVICE=user-service ACTION=up"
