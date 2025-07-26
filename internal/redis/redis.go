package redis

import (
	"context"
	"crypto/tls"
	"fmt"

	"github/erastusk/tracer/internal/types"
	"github/erastusk/tracer/internal/utils"

	"github.com/redis/go-redis/v9"
)

// func Redis(addr, user, pass string, sslskip bool) *redis.StatusCmd {
type redisApp struct {
	SslSkip bool
}

func NewRedis() *redisApp {
	return &redisApp{}
}

func (r *redisApp) Connectivity(opts types.PromptOptions) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     opts.Endpoint,
		Username: opts.Username,
		Password: opts.Password,
		DB:       0,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: r.SslSkip,
		},
	}).Ping(context.Background())
	if rdb.Err() != nil {
		utils.Red.Println(fmt.Sprintf("Redis connectivity to :%s was unsuccessful!", opts.Endpoint))
	} else {
		utils.Green.Println("Redis connectivity was successful!")
	}
}
