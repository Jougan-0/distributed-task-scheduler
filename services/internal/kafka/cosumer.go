package kafka

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type ConsumerGroupHandler struct{}

func (ConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (ConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		value := string(message.Value)

		parts := strings.Split(value, ":")
		if len(parts) != 2 {
			log.Printf("Kafka Consumer: Invalid message format: %s", value)
			continue
		}

		taskID, err := uuid.Parse(strings.TrimSpace(parts[1]))
		if err != nil {
			log.Printf("Kafka Consumer: Invalid UUID in message: %s", parts[1])
			continue
		}
		AddEvent(value)
		log.Printf("Kafka Consumer: Received event '%s' for Task ID=%s", parts[0], taskID.String())

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
