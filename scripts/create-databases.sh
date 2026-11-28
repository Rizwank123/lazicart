#!/usr/bin/env bash
set -euo pipefail

# --------------------------------------------------
# Create all databases for LahziCart services
# --------------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "📂 Project root: $ROOT_DIR"

# Database configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-root}"

PGPASSWORD="$DB_PASSWORD"
export PGPASSWORD

echo "🗄️  Checking/creating databases..."

# Define all databases needed by services and their display names
declare -A DATABASES=(
    ["user_auth"]="Auth Service"
    ["user-service"]="User Service"
    ["catalog"]="Catalog Service"
    ["inventory"]="Inventory Service"
    ["cart"]="Cart Service"
    ["orders"]="Order Service"
    ["payment"]="Payment Service"
    ["notification"]="Notification Service"
    ["delivery"]="Delivery Service"
)

CREATE_FAILED=0
CREATE_SKIPPED=0
CREATE_NEW=0

for DB_NAME in "${!DATABASES[@]}"; do
    SERVICE_NAME="${DATABASES[$DB_NAME]}"

    # Check if database already exists using psql -lqt
    # The quoting approach: quote the name for the grep pattern
    if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -lqt 2>/dev/null | \
       cut -d '|' -f 1 | \
       grep -qw "$DB_NAME"; then
        echo "✅ Database '$DB_NAME' already exists (for $SERVICE_NAME)"
        CREATE_SKIPPED=$((CREATE_SKIPPED + 1))
        continue
    fi

    # Create the database (quote identifier to handle names with hyphens/ special chars)
    echo "🔧 Creating database '$DB_NAME' for $SERVICE_NAME..."
    if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -c "CREATE DATABASE \"$DB_NAME\";" 2>/dev/null; then
        echo "✅ Database '$DB_NAME' created successfully"
        CREATE_NEW=$((CREATE_NEW + 1))
    else
        echo "❌ Failed to create database '$DB_NAME'"
        CREATE_FAILED=$((CREATE_FAILED + 1))
    fi
done

echo ""
echo "========================================="
echo "Creation Summary"
echo "========================================="
echo "✅ Created new: $CREATE_NEW"
echo "   Already exist: $CREATE_SKIPPED"
echo "❌ Failed: $CREATE_FAILED"
echo "========================================="
echo "Total: $((CREATE_NEW + CREATE_SKIPPED + CREATE_FAILED)) databases processed"

if [ "$CREATE_FAILED" -gt 0 ]; then
    echo "Some databases could not be created. Check the errors above."
    exit 1
fi

echo "All databases are ready!"