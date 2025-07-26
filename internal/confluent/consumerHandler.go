package confluent

import (
	"fmt"

	"github.com/IBM/sarama"
)

type consumerGroupHandler struct{}

func (consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fmt.Printf("Consumed: topic=%v partition=%v offset=%v key=%s value=%s\n",
			message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
		session.MarkMessage(message, "")
	}
	return nil
}
