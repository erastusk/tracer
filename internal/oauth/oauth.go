package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github/erastusk/tracer/internal/types"

	"github.com/IBM/sarama"
)

type OAuthTokenProvider struct {
	clientID       string
	clientSecret   string
	tokenURL       string
	scopes         string
	logicalCluster string
	poolID         string
}

func NewOAuthTokenProvider(s types.Secrets) *OAuthTokenProvider {
	return &OAuthTokenProvider{
		clientID:       s.ConfluentKafkaOauthClientID,
		clientSecret:   s.ConfluentKafkaOauthClientSecret,
		tokenURL:       s.TokenURL,
		scopes:         s.Scope,
		logicalCluster: s.ConfluentKafkaOauthLogicalCluster,
		poolID:         s.ConfluentKafkaOauthIdentityPoolID,
	}
}

func (o *OAuthTokenProvider) Token() (*sarama.AccessToken, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", o.clientID)
	data.Set("client_secret", o.clientSecret)
	data.Set("scope", o.scopes)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", o.tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("req err", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("do err", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp types.OauthTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		fmt.Println("decode err", err)
		return nil, err
	}

	if tokenResp.AccessToken == "" {
		return nil, errors.New("missing access_token in response")
	}

	return &sarama.AccessToken{
		Token: tokenResp.AccessToken,
		Extensions: map[string]string{
			"logicalCluster": o.logicalCluster,
			"identityPoolId": o.poolID,
		},
	}, nil
}

func (o *OAuthTokenProvider) NewConfigOauthMap() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Net.SASL.Enable = true
	config.Producer.Return.Successes = true
	config.Net.SASL.Mechanism = sarama.SASLTypeOAuth
	config.Net.SASL.TokenProvider = o
	config.Net.TLS.Enable = true
	return config
}
