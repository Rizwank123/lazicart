#!/usr/bin/env bash
set -euo pipefail

# --------------------------------------------------
# Resolve paths safely
# --------------------------------------------------
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
SERVICE_NAME="delivery-service"
SERVICE_DIR="$ROOT_DIR/Services/$SERVICE_NAME"

cd "$ROOT_DIR"

echo "📂 Project root : $ROOT_DIR"
echo "📦 Service root : $SERVICE_DIR"

# --------------------------------------------------
# Load service-level .env
# --------------------------------------------------
ENV_FILE="$SERVICE_DIR/.env"

if [ -f "$ENV_FILE" ]; then
	echo "🔐 Loading env from $ENV_FILE"
	set -a
	. "$ENV_FILE"
	set +a
else
	echo "❌ .env not found at $ENV_FILE"
	exit 1
fi

# --------------------------------------------------
# Validate DB variables
# --------------------------------------------------
REQUIRED_VARS=(DB_USER DB_PASSWORD DB_HOST DB_PORT DB_NAME)

for var in "${REQUIRED_VARS[@]}"; do
	if [ -z "${!var:-}" ]; then
		echo "❌ Missing environment variable: $var"
		exit 1
	fi
done

# --------------------------------------------------
# Migration directory
# --------------------------------------------------
MIGRATIONS_DIR="$SERVICE_DIR/internal/db/migration"

if [ ! -d "$MIGRATIONS_DIR" ]; then
	echo "❌ Migration directory not found: $MIGRATIONS_DIR"
	exit 1
fi

# --------------------------------------------------
# Action (up / down)
# --------------------------------------------------
ACTION="${1:-}"

if [[ "$ACTION" != "up" && "$ACTION" != "down" ]]; then
	echo "❌ Invalid action: $ACTION"
	echo "✅ Usage: up | down"
	exit 1
fi

# --------------------------------------------------
# Build DB DSN
# --------------------------------------------------
DB_DSN="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

echo "🗄️  DB         : postgres://${DB_USER}@${DB_HOST}:${DB_PORT}/${DB_NAME}"
echo "📦 Migrations : $MIGRATIONS_DIR"
echo "⚙️  Action     : $ACTION"

# --------------------------------------------------
# Run goose
# --------------------------------------------------
goose -dir "$MIGRATIONS_DIR" postgres "$DB_DSN" "$ACTION"
