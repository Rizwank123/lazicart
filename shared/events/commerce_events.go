package events

import "time"

type OrderItemEvent struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price,omitempty"`
}

type OrderCreatedEvent struct {
	EventType string           `json:"event_type"`
	Version   int              `json:"version"`
	OrderID   string           `json:"order_id"`
	UserID    string           `json:"user_id"`
	Items     []OrderItemEvent `json:"items"`
	Total     float64          `json:"total"`
	CreatedAt time.Time        `json:"created_at"`
}

type OrderCancelledEvent struct {
	EventType   string    `json:"event_type"`
	Version     int       `json:"version"`
	OrderID     string    `json:"order_id"`
	UserID      string    `json:"user_id"`
	CancelledAt time.Time `json:"cancelled_at"`
}

type InventoryReservedEvent struct {
	EventType  string           `json:"event_type"`
	Version    int              `json:"version"`
	OrderID    string           `json:"order_id"`
	Items      []OrderItemEvent `json:"items"`
	ReservedAt time.Time        `json:"reserved_at"`
}

type InventoryReleasedEvent struct {
	EventType  string           `json:"event_type"`
	Version    int              `json:"version"`
	OrderID    string           `json:"order_id"`
	Items      []OrderItemEvent `json:"items"`
	ReleasedAt time.Time        `json:"released_at"`
}

type PaymentEvent struct {
	EventType     string    `json:"event_type"`
	Version       int       `json:"version"`
	PaymentID     string    `json:"payment_id"`
	OrderID       string    `json:"order_id"`
	UserID        string    `json:"user_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	TransactionID string    `json:"transaction_id,omitempty"`
	OccurredAt    time.Time `json:"occurred_at"`
}
