package confluent

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github/erastusk/tracer/internal/oauth"
	"github/erastusk/tracer/internal/secrets"
	"github/erastusk/tracer/internal/types"
	"github/erastusk/tracer/internal/utils"

	"github.com/IBM/sarama"
)

var (
	secretPath = "/dev/ftb/sre/secrets"
	region     = "us-east-1"
)

//go:generate mockgen -destination=../../mocks/confluentoauth.go -package=mocks github/erastusk/tracer/internal/confluent ConfluentOauth
type ConfluentOauth interface {
	types.KafkaApps
}

type ConfluentOauthImpl struct {
	producerClient sarama.SyncProducer
	consumerClient sarama.Consumer
	client         sarama.Client
	tokenProvider  *oauth.OAuthTokenProvider
	accessToken    *sarama.AccessToken
	secrets        types.Secrets
	partitions     []int32
	topic          string
	error          error
	config         *sarama.Config
	handler        consumerGroupHandler
}

func NewConfluentOauth() (*ConfluentOauthImpl, error) {
	s := secrets.SecretsManagerImpl()
	sess, err := s.GetSession(secretPath, region)
	if err != nil {
		return nil, fmt.Errorf("unable to create session: %w", err)
	}
	secret, err := sess.GetSecrets()
	if err != nil {
		return nil, fmt.Errorf("unable to get secrets: %w", err)
	}
	c := oauth.NewOAuthTokenProvider(secret)
	config := c.NewConfigOauthMap()
	client, err := sarama.NewClient([]string{secret.ConfluentKafkaServer}, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create OauthToken: %w", err)
	}
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, fmt.Errorf("newsync producer err: %w", err)
	}
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, fmt.Errorf("newsync producer err: %w", err)
	}
	return &ConfluentOauthImpl{
		secrets:        secret,
		tokenProvider:  c,
		client:         client,
		producerClient: producer,
		consumerClient: consumer,
		config:         config,
		handler:        consumerGroupHandler{},
	}, nil
}

func (o *ConfluentOauthImpl) Connectivity(types.PromptOptions) error {
	if o.client == nil {
		return fmt.Errorf("client failed to connect")
	}
	return nil
}

func (o *ConfluentOauthImpl) Produce(p types.PromptOptions, topic string) error {
	defer o.producerClient.Close()
	payload := types.UserEvent{
		UserID:    rand.Intn(100),
		Event:     "Confluent Oauth Tracer logs",
		Timestamp: time.Now(),
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(payloadBytes),
		Partition: int32(0),
	}
	partition, offset, err := o.producerClient.SendMessage(msg)
	if err != nil || partition < 0 || offset < 0 {
		return err
	} else {
		utils.Green.Printf("**********************************\nDelivered to partition %d at offset %d topic %s\n**********************************\n",
			partition, offset, topic)
	}
	return nil
}

func (o *ConfluentOauthImpl) ListTopics(p types.PromptOptions) error {
	defer o.client.Close()
	topics, err := o.client.Topics()
	if err != nil {
		return err
	}
	for _, topic := range topics {
		utils.Green.Println("topic: ", topic)
	}
	return nil
}

func (o *ConfluentOauthImpl) Consume(p types.PromptOptions, topic string) error {
	cg, err := sarama.NewConsumerGroup([]string{o.secrets.ConfluentKafkaServer}, "eisl-demo", o.config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer cg.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	for ctx.Err() == nil {
		err := cg.Consume(ctx, []string{topic}, o.handler)
		if err != nil {
			o.error = err
		}
	}
	return o.error
}
