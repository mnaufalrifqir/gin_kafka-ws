package kafka

import (
	"go_confluence/websocket"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func SetupKafkaProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Timeout = 5 * time.Second
	config.Producer.Return.Successes = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

func SendUserToKafka(producer sarama.SyncProducer, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic:     "my-topic",
		Value:     sarama.StringEncoder(message),
		Timestamp: time.Now(),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Error sending message to Kafka: %v", err)
	}

	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
	return err
}

func StartKafkaConsumer() {
	// Create Kafka consumer
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	topic := "my-topic"
	kafkaHost := "localhost:9092"

	consumer, err := sarama.NewConsumer([]string{kafkaHost}, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %s", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Failed to create partition consumer: %s", err)
	}
	defer partitionConsumer.Close()

	HandleMessage(partitionConsumer)

	// group, err := sarama.NewConsumerGroupFromClient(kafkaGroupID, consumer)
	// if err != nil {
	// 	return err
	// }
	// defer group.Close()

	// go func() {
	// 	for err := range group.Errors() {
	// 		log.Printf("Consumer group error: %v", err)
	// 	}
	// }()
	// // Start consuming messages

	// handler := &ConsumerGroupHandler{}
	// wg := &sync.WaitGroup{}
	// wg.Add(1)

	// ctx := context.Background()

	// go func() {
	// 	defer wg.Done()
	// 	for {
	// 		err := group.Consume(ctx, []string{topic}, handler)
	// 		if err != nil {
	// 			log.Printf("Error from consumer: %v", err)
	// 		}
	// 		// Exit the loop if the consumer group is closed or encounters an error
	// 		if handler.Closed() {
	// 			return
	// 		}
	// 	}
	// }()

	// wg.Wait()
	// return nil
}

func HandleMessage(partitionConsumer sarama.PartitionConsumer) {
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			websocket.SendWebSocketMessage(string(msg.Value))
		case err := <-partitionConsumer.Errors():
			log.Printf("Error: %v\n", err.Err)
		}
	}
}
