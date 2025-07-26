package types

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type PromptOptions struct {
	Endpoint string
	Username string
	Password string
}

type (
	Config    *kafka.ConfigMap
	UserEvent struct {
		UserID    int
		Event     string
		Timestamp time.Time
	}
)

type Secrets struct {
	KafkaUser                         string `json:"kafka_user"`
	KafkaPassword                     string `json:"kafka_password"`
	DBUser                            string `json:"db_user"`
	DBPassword                        string `json:"db_password"`
	Username                          string `json:"username"`
	Password                          string `json:"password"`
	ConfluentKafkaOauthClientID       string `json:"confluent_kafka_oauth_client_id"`
	ConfluentKafkaOauthClientSecret   string `json:"confluent_kafka_oauth_client_secret"`
	ConfluentKafkaOauthLogicalCluster string `json:"confluent_kafka_oauth_logical_cluster"`
	ConfluentKafkaOauthIdentityPoolID string `json:"confluent_kafka_oauth_identity_pool_id"`
	TokenURL                          string `json:"token_url"`
	APIKey                            string `json:"api_key"`
	APISecret                         string `json:"api_secret"`
	ConfluentKafkaServer              string `json:"confluent_kafka_server"`
	Scope                             string `json:"scope"`
	RedisUsername                     string `json:"redis_username"`
	RedisPassword                     string `json:"redis_password"`
	RedisServer                       string `json:"redis_server"`
}

type OauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}
