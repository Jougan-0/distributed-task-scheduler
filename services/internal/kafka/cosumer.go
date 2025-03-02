package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type KafkaEvent struct {
	Event         string    `json:"event"`
	TaskID        uuid.UUID `json:"task_id"`
	TaskName      string    `json:"task_name"`
	TaskType      string    `json:"task_type"`
	Priority      int       `json:"priority"`
	Attempts      int       `json:"attempts,omitempty"`
	MaxRetries    int       `json:"max_retries,omitempty"`
	ScheduledTime string    `json:"scheduled_time"`
	CompletedAt   string    `json:"completed_at,omitempty"`
	FailedAt      string    `json:"failed_at,omitempty"`
}
type ConsumerGroupHandler struct{}

func (ConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (ConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var event KafkaEvent

		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("Kafka Consumer: Failed to parse JSON: %v, Message: %s", err, string(message.Value))
			continue
		}

		AddEvent(event)
		log.Printf("Kafka Consumer: Received event '%s' for Task ID=%s", event.Event, event.TaskID.String())

		session.MarkMessage(message, "")
	}
	return nil
}

func StartConsumerGroup(ctx context.Context, groupID, topic string) error {
	brokersEnv := os.Getenv("KAFKA_BROKERS")
	if brokersEnv == "" {
		brokersEnv = "localhost:9092"
	}
	brokers := strings.Split(brokersEnv, ",")

	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return err
	}
	defer consumerGroup.Close()

	handler := ConsumerGroupHandler{}

	for {
		if err := consumerGroup.Consume(ctx, []string{topic}, handler); err != nil {
			log.Printf("Kafka Consumer: Error during consumption: %v", err)
			time.Sleep(5 * time.Second)
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}
