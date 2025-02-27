package kafka

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

var Producer sarama.SyncProducer

func InitProducer() error {
	brokersEnv := os.Getenv("KAFKA_BROKERS")
	if brokersEnv == "" {
		brokersEnv = "localhost:9092"
	}
	brokers := strings.Split(brokersEnv, ",")

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 10 * time.Second
	config.Producer.Compression = sarama.CompressionSnappy

	var err error
	Producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return err
	}
	log.Println("Kafka Producer initialized.")
	return nil
}

func PublishMessage(topic, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := Producer.SendMessage(msg)
	if err != nil {
		log.Printf("Kafka Producer: Failed to send message: %v", err)
		return err
	}
	log.Printf("Kafka Producer: Message sent to topic %s partition %d offset %d", topic, partition, offset)
	return nil
}
