package observability

import "github.com/prometheus/client_golang/prometheus"

var (
	KafkaMessagesConsumed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_messages_consumed_total",
			Help: "Total number of messages consumed from Kafka",
		},
		[]string{"topic", "status"},
	)

	KafkaProcessingDration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "kafka_processing_duration_seconds",
			Help: "Duration of processing Kafka messages",
		},
		[]string{"topic"},
	)
)

func Register() {
	prometheus.MustRegister(
		KafkaMessagesConsumed,
		KafkaProcessingDration,
	)
}
