package events

import (
	"fmt"
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

func EnsureTopic(brokers []string, topic string, partitions int, replicationFactor int) error {
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Second))

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s", "localhost:9092"))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	return controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	})
}
