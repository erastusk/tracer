package secrets

import (
	"context"
	"encoding/json"
	"log"

	"github/erastusk/tracer/internal/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

//go:generate mockgen -destination=../../mocks/secretsmanager.go -package=mocks github/erastusk/tracer/internal/secrets SecretsManager
type SecretsManager interface {
	GetSession(name string, region string) (*SecretsSession, error)
	GetSecrets() (types.Secrets, error)
}
type (
	SecretsSession struct {
		SecretKey string
		Region    string
		Sess      *secretsmanager.Client
	}
)

func SecretsManagerImpl() SecretsManager {
	return &SecretsSession{}
}

func (s *SecretsSession) GetSession(name string, region string) (*SecretsSession, error) {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	return &SecretsSession{
		SecretKey: name, Region: region,
		Sess: secretsmanager.NewFromConfig(config),
	}, nil
}

func (s *SecretsSession) GetSecrets() (types.Secrets, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(s.SecretKey),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}
	result, err := s.Sess.GetSecretValue(context.TODO(), input)
	if err != nil {
		return types.Secrets{}, err
	}
	var sec types.Secrets
	if err := json.Unmarshal([]byte(*result.SecretString), &sec); err != nil {
		return types.Secrets{}, err
	}
	return sec, nil
}
