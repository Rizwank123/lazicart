#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

SERVICES=(
    "auth-service"
    "user-service"
    "catalog-service"
    "inventory-service"
    "cart-service"
    "order-service"
    "payment-service"
    "notification-service"
    "delivery-service"
    "api-gateway"
)

PIDS=()

cleanup() {
    echo ""
    echo "Stopping services..."
    for pid in "${PIDS[@]:-}"; do
        if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
            kill "$pid" 2>/dev/null || true
        fi
    done
    wait 2>/dev/null || true
    echo "All services stopped."
}

trap cleanup INT TERM EXIT

echo "Project root: $ROOT_DIR"
echo "Starting LahziCart services..."

for service in "${SERVICES[@]}"; do
    SERVICE_DIR="$ROOT_DIR/Services/$service"

    if [ ! -d "$SERVICE_DIR" ]; then
        echo "Service directory not found: $SERVICE_DIR"
        exit 1
    fi

    if [ ! -d "$SERVICE_DIR/cmd/server" ]; then
        echo "Server entrypoint not found for $service: $SERVICE_DIR/cmd/server"
        exit 1
    fi

    echo "Starting $service..."
    (
        cd "$SERVICE_DIR"
        GOTOOLCHAIN="${GOTOOLCHAIN:-auto}" go run -p 2 ./cmd/server
    ) &
    PIDS+=("$!")

    # Delay to stagger startup and manage system memory
    sleep 5
done

echo ""
echo "Started ${#PIDS[@]} services."
echo "Press Ctrl+C to stop all services."

wait
