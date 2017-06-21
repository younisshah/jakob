package jkafka

import (
	"log"
	"os"

	"os/signal"

	"github.com/Shopify/sarama"
)

/**
*  Created by Galileo on 21/6/17.
 */

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
func Consume() error {

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

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	done := make(chan bool)

	go func() {
		for {
			select {
			case err := <-partitionConsumer.Errors():
				logger.Println("error while consuming ", err)
			case cmd := <-partitionConsumer.Messages():
				logger.Printf("rcvd >> cmd: %s", string(cmd.Value))
			case <-shutdown:
				logger.Println("kafka consumer shutdown triggered")
				done <- true
			}
		}
	}()
	<-done
	logger.Println("Exit.")
	return nil
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
