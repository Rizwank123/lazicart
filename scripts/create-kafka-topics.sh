#!/usr/bin/env bash
# Create the Kafka topics the services publish to.
#
# The topics are not created automatically by the broker in this setup, so a
# publish to a missing topic fails. That turns a successful database write
# into a failed API response, for example when the user-service publishes
# user.updated after a profile update.
#
# The list is kept in sync with shared/events/topics.go.
set -euo pipefail

BOOTSTRAP_SERVER="${KAFKA_BOOTSTRAP_SERVER:-localhost:9092}"

TOPICS=(
  user.created
  user.updated
  user.deleted
  order.created
  order.cancelled
  payment.succeeded
  payment.failed
  notification.created
  inventory.reserved
  inventory.released
  inventory.adjusted
)

echo "Creating Kafka topics on ${BOOTSTRAP_SERVER}..."
for topic in "${TOPICS[@]}"; do
  if kafka-topics --bootstrap-server "$BOOTSTRAP_SERVER" --list 2>/dev/null | grep -qx "$topic"; then
    echo "  = $topic (already exists)"
    continue
  fi
  kafka-topics --bootstrap-server "$BOOTSTRAP_SERVER" --create --topic "$topic" --partitions 1 --replication-factor 1 >/dev/null
  echo "  + $topic"
done

echo "Done."
