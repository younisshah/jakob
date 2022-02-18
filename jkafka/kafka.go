package jkafka

import (
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/younisshah/jakob/replicate"
)

const (
	_KAFKA_TOPIC     = "jakob_topic"
	_KAFKA_CLIENT_ID = "jakob_client"
)

var (
	hosts  = []string{"localhost:9092"}
	logger = log.New(os.Stderr, "[jakob-kafka] ", log.LstdFlags)
)

// Consume consumes from a Kafka topic
// Use this to sync commands to a new getter peer
func Consume(peer string) error {

	consumer, err := sarama.NewConsumer(hosts, nil)

	if err != nil {
		logger.Println("couldn't create Kafka consumer. Hosts are: ", hosts)
		logger.Println(" -ERROR", err)
		return err
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			logger.Println("couldn't close kafka consumer")
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(_KAFKA_TOPIC, 0, sarama.OffsetOldest)
	if err != nil {
		logger.Println("couldn't creates a PartitionConsumer on the given topic/partition.")
		logger.Println(" -TOPIC", _KAFKA_TOPIC)
		logger.Println(" -PARTITION", 0)
		return err
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			logger.Println("couldn't close PartitionConsumer")
		}
	}()

	var i int64
	messages := make(map[int64]map[string][]interface{})
DONE:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			value := string(msg.Value)
			tokens := strings.Fields(value)
			cmdName := tokens[0]
			args := pack(tokens[1:])
			cmd := map[string][]interface{}{cmdName: args}
			messages[i] = cmd
			i = i + 1
			if i == partitionConsumer.HighWaterMarkOffset() {
				break DONE
			}
		}
	}

	err = replicate.Replicate(peer, messages)

	logger.Println("Exit.")
	return err
}

// Produce produces the cmd to a kafka topic
func Produce(cmd string) error {
	config := getConfig()
	producer, err := sarama.NewSyncProducer(hosts, config)
	if err != nil {
		logger.Println("couldn't create kafka producer", err)
		return err
	}
	defer func() {
		if err := producer.Close(); err != nil {
			logger.Println("couldn't close producer", err)
		}
	}()
	msg := &sarama.ProducerMessage{
		Topic: _KAFKA_TOPIC,
		Value: sarama.StringEncoder(cmd),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		logger.Printf("couldn't produce cmd [%s] on topic [%s]", cmd, _KAFKA_TOPIC)
		logger.Println(" -ERROR", err)
		return err
	}
	logger.Printf("command stored in topic [%s], partition [%d], offest [%d]\n", _KAFKA_TOPIC, partition, offset)
	return nil
}

func getConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5
	config.ClientID = _KAFKA_CLIENT_ID
	return config
}

func pack(s []string) []interface{} {
	args := make([]interface{}, len(s))
	for i := range s {
		args[i] = s[i]
	}
	return args
}
