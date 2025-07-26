package confluent

import (
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	"github/erastusk/tracer/internal/types"
	"github/erastusk/tracer/internal/utils"

	"github.com/IBM/sarama"
)

type confluent struct {
	producerClient sarama.SyncProducer
	consumerClient sarama.Consumer
	partitions     []int32
	topic          string
	error          error
}

func NewConfluent() *confluent {
	return &confluent{}
}

func (c *confluent) Connectivity(opts types.PromptOptions) error {
	broker := []string{opts.Endpoint}
	config := NewConfigMap(opts.Username, opts.Password)
	a, err := sarama.NewClient(broker, config)
	if err != nil {
		return err
	}
	defer a.Close()
	return nil
}

func (c *confluent) Produce(opts types.PromptOptions, topic string) error {
	c.topic = topic
	c.newClient(opts, "producer")
	if c.error != nil {
		// utils.Red.Println("Unable to create producer", c.error)
		return c.error
	}
	payload := types.UserEvent{
		UserID:    rand.Intn(100),
		Event:     "Tracer logs",
		Timestamp: time.Now(),
	}
	wg := &sync.WaitGroup{}
	produceMessage := func() {
		go func() {
			defer wg.Done()
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				log.Fatal(err)
			}
			defer c.producerClient.Close()
			msg := &sarama.ProducerMessage{
				Topic:     topic,
				Value:     sarama.ByteEncoder(payloadBytes),
				Partition: int32(0),
			}
			partition, offset, err := c.producerClient.SendMessage(msg)
			if err != nil || partition < 0 || offset < 0 {
				c.error = err
			} else {
				utils.Green.Printf("**********************************\nDelivered to partition %d at offset %d topic %s\n**********************************\n",
					partition, offset, topic)
			}
		}()
	}
	wg.Add(1)
	produceMessage()
	wg.Wait()
	if c.error != nil {
		return c.error
	}
	return nil
}

func (c *confluent) Consume(opts types.PromptOptions, topic string) error {
	c.topic = topic
	c.newClient(opts, "consumer")
	if c.error != nil {
		// utils.Red.Println("Unable to create consumer", c.error)
		return c.error
	}
	readMessages := func(done chan any) error {
		go func() {
			defer c.consumerClient.Close()
			msg, err := c.consumerClient.ConsumePartition(topic, int32(0), sarama.OffsetOldest)
			if err != nil {
				// utils.Red.Println("Unable to start partition consumer")
				c.error = err
			}
			select {
			case m := <-msg.Messages():
				utils.Green.Printf("**********************************\npartition: %d\noffset: %d\ntopic: %s\nmessage: %s\n**********************************\n",
					m.Partition, m.Offset, c.topic, string(m.Value))
			case <-done:
			}
		}()
		return c.error
	}
	done := make(chan any)
	c.error = readMessages(done)
	if c.error != nil {
		return c.error
	}
	time.Sleep(5 * time.Second)
	close(done)
	return nil
}

func (c *confluent) newClient(opts types.PromptOptions, action string) {
	servers := opts.Endpoint
	username := opts.Username
	password := opts.Password
	switch action {
	case "clientAdmin":
		c.newClientAdmin(servers, username, password)
	case "producer":
		c.newProducer(servers, username, password)
	case "consumer":
		c.newConsumer(servers, username, password)
	}
}

func NewConfigMapOauth(ccloudAPIKey, ccloudAPISecret string) *sarama.Config {
	return nil
}

func NewConfigMap(ccloudAPIKey, ccloudAPISecret string) *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	// SASL Authentication configuration
	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext // or sarama.SASLTypeSCRAMSHA256, etc.
	config.Net.SASL.User = ccloudAPIKey
	config.Net.SASL.Password = ccloudAPISecret
	config.Version = sarama.V3_0_0_0
	config.Net.DialTimeout = time.Second * 1
	config.Metadata.Timeout = time.Second * 1
	config.Metadata.Retry.Backoff = 0
	config.Metadata.Retry.Max = 0

	// Enable TLS if required (commonly needed with SASL)
	config.Net.TLS.Enable = true
	return config
}

func (c *confluent) newClientAdmin(servers, ccloudAPIKey, ccloudAPISecret string) {}
func (c *confluent) newProducer(servers, ccloudAPIKey, ccloudAPISecret string) error {
	broker := []string{servers}
	config := NewConfigMap(ccloudAPIKey, ccloudAPISecret)
	a, err := sarama.NewSyncProducer(broker, config)
	if err != nil {
		c.error = err
		return err
	}
	c.producerClient = a
	return nil
}

func (c *confluent) newConsumer(servers, ccloudAPIKey, ccloudAPISecret string) {
	broker := []string{servers}
	config := NewConfigMap(ccloudAPIKey, ccloudAPISecret)
	a, err := sarama.NewConsumer(broker, config)
	if err != nil {
		c.error = err
		return
	}
	part, err := a.Partitions(c.topic)
	if err != nil {
		c.error = err
		return
	}
	c.partitions = part
	c.consumerClient = a
}

func (c *confluent) ListTopics(opts types.PromptOptions) error {
	broker := []string{opts.Endpoint}
	config := NewConfigMap(opts.Username, opts.Password)
	a, err := sarama.NewClient(broker, config)
	if err != nil {
		c.error = err
		return err
	}
	defer a.Close()
	topics, err := a.Topics()
	if err != nil {
		c.error = err
		return err
	}
	for _, topic := range topics {
		utils.Green.Println("topic: ", topic)
	}
	return nil
}
