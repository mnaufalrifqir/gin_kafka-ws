package kafka

import (
	"fmt"
	websocketGo "go_confluence/websocket"
	"os"
	"time"

	"github.com/IBM/sarama"
)

type ConsumerGroupHandler struct {
	closed bool
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// Process the consumed message as desired

		// Creating timestamp
		start := message.Timestamp

		str := string(message.Value) // Convert byte array to string
		fmt.Println(str)             // print out the message as a string
		websocketGo.SendWebSocketMessage(str)

		processing := time.Since(start)
		fmt.Fprintf(os.Stdout, "\033[0;31m Time taken: %s\033[0m\n ", processing)
		session.MarkMessage(message, "consumed") //MarkMessage marks a message as consumed.

	}
	return nil
}

func (h *ConsumerGroupHandler) Closed() bool {
	return h.closed
}
