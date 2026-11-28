package events

import (
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	UserCreated         = "user.created"
	OrderCreated        = "order.created"
	OrderCancelled      = "order.cancelled"
	PaymentSucceeded    = "payment.succeeded"
	PaymentFailed       = "payment.failed"
	NotificationCreated = "notification.created"
	UserUpdated         = "user.updated"
	UserDeleted         = "user.deleted"
	InventoryReserved   = "inventory.reserved"
	InventoryReleased   = "inventory.released"
	InventoryAdjusted   = "inventory.adjusted"
)

// AllTopics lists every topic the services publish to.
var AllTopics = []string{
	UserCreated,
	UserUpdated,
	UserDeleted,
	OrderCreated,
	OrderCancelled,
	PaymentSucceeded,
	PaymentFailed,
	NotificationCreated,
	InventoryReserved,
	InventoryReleased,
	InventoryAdjusted,
}

// EnsureTopic creates a topic when it does not exist yet.
//
// The broker does not create topics automatically in this setup, so a publish
// to a missing topic fails with "Unknown Topic Or Partition". A topic that
// already exists is not an error.
func EnsureTopic(brokers []string, topic string, partitions int, replicationFactor int) error {
	if len(brokers) == 0 {
		return errors.New("no kafka brokers configured")
	}

	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Second))

	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	})
	if err != nil && !errors.Is(err, kafka.TopicAlreadyExists) {
		return err
	}
	return nil
}
