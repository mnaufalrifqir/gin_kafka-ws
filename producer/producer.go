package main

import (
	"encoding/json"
	"log"
	"net/http"

	"go_confluence/kafka"

	"github.com/gin-gonic/gin"
)

type MessageTest struct {
	MessageTest string `json:"message"`
}

func main() {
	router := gin.Default()

	producer, err := kafka.SetupKafkaProducer()
	if err != nil {
		log.Println("Failed to setup kafka:", err)
		return
	}
	defer producer.Close()

	router.GET("/send-data", func(c *gin.Context) {
		var messageTest MessageTest

		// // Bind JSON request body to the user struct
		if err := c.ShouldBindJSON(&messageTest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		jsonData, err := json.Marshal(messageTest)
		if err != nil {
			log.Fatalf("Error encoding JSON: %v", err)
		} else {
			err := kafka.SendUserToKafka(producer, jsonData)
			if err != nil {
				log.Print("Failed to produce message:")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to produce message"})
				return
			}
		}

		c.JSON(http.StatusOK, messageTest)
	})

	// Start the Gin server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start servers: %s", err)
	}

	// err = kafka.SendUserToKafka(producer, []byte("Hello, Kafka!"))
	// if err != nil {
	// 	log.Print("Failed to produce message:")
	// 	return
	// }

	// config := sarama.NewConfig()
	// config.Producer.Return.Successes = true

	// producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	// if err != nil {
	// 	log.Fatalf("Error creating producer: %v", err)
	// }
	// defer producer.Close()

	// message := &sarama.ProducerMessage{
	// 	Topic: "my-topic",
	// 	Value: sarama.StringEncoder("Hello, Kafka!"),
	// }

	// // Kirim pesan ke Kafka
	// partition, offset, err := producer.SendMessage(message)
	// if err != nil {
	// 	log.Fatalf("Error sending message to Kafka: %v", err)
	// }

	// log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
}
