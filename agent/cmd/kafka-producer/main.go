package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	var (
		brokers  = flag.String("brokers", "127.0.0.1:19093", "comma-separated Kafka broker addresses")
		topic    = flag.String("topic", "hotgo-test-events", "Kafka topic")
		key      = flag.String("key", "", "optional Kafka message key")
		message  = flag.String("message", "", "message payload; generated JSON is used when empty")
		count    = flag.Int("count", 100000, "number of messages to send")
		username = flag.String("username", "hotgo", "SASL username")
		password = flag.String("password", "hotgo-secret", "SASL password")
	)
	flag.Parse()

	if *topic == "" {
		log.Fatal("topic is required")
	}
	if *count <= 0 {
		log.Fatal("count must be greater than 0")
	}

	producer, err := newProducer(splitCSV(*brokers), *username, *password)
	if err != nil {
		log.Fatalf("create kafka producer failed: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("close kafka producer failed: %v", err)
		}
	}()

	for i := 1; i <= *count; i++ {
		payload := *message
		if payload == "" {
			payload = fmt.Sprintf(`{"source":"hotgo-kafka-producer","seq":%d,"time":"%s","message":"hello kafka"}`, i, time.Now().Format(time.RFC3339Nano))
		}

		msg := &sarama.ProducerMessage{
			Topic: *topic,
			Value: sarama.StringEncoder(payload),
		}
		if *key != "" {
			msg.Key = sarama.StringEncoder(*key)
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Fatalf("send message %d failed: %v", i, err)
		}
		fmt.Fprintf(os.Stdout, "sent topic=%s partition=%d offset=%d payload=%s\n", *topic, partition, offset, payload)
		time.Sleep(time.Second)
	}
}

func newProducer(brokers []string, username, password string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V3_5_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	config.Net.DialTimeout = 10 * time.Second
	config.Net.ReadTimeout = 30 * time.Second
	config.Net.WriteTimeout = 30 * time.Second

	if username != "" || password != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.User = username
		config.Net.SASL.Password = password
	}

	return sarama.NewSyncProducer(brokers, config)
}

func splitCSV(value string) []string {
	var values []string
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			values = append(values, item)
		}
	}
	return values
}
