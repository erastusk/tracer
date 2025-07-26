package Cmd

import (
	"github/erastusk/tracer/internal/cache"
	"github/erastusk/tracer/internal/prompts"
	"github/erastusk/tracer/internal/redis"
	"github/erastusk/tracer/internal/types"
	"github/erastusk/tracer/internal/utils"

	"github.com/spf13/cobra"
)

type RedisType struct {
	Type  string
	Cache cache.Cache
}

func NewRedisType(redisType string) *RedisType {
	cache := cache.NewCache(redisType)
	return &RedisType{
		Type:  redisType,
		Cache: *cache,
	}
}

// redisCmd represents the redis command
func (r *RedisType) Cmd() *cobra.Command {
	// var msk RedisLastAnswers
	// Load the last answers from the cach
	var c types.PromptOptions
	b := redis.NewRedis()
	count := 0
	// EncodeAnyToStruct(c, &msk)
	return &cobra.Command{
		Use:   "redis",
		Short: "Redis connectivity check",
		Run: func(cmd *cobra.Command, args []string) {
			options := []string{"Redis connectivity test"}
			for {
				count++
				if count > 0 {
					c = r.Cache.LoadCache()
				}
				select {
				case <-utils.Done:
					utils.Green.Println("exiting....")
					return
				default:
					switch prompts.GetUserPrompt(options) {
					case "Redis connectivity test":
						z, err := prompts.GetPrompts(c, utils.Options)
						if err != nil {
							utils.Red.Println(err)
						} else {
							r.Cache.SaveCache(z)
							s := prompts.GetUserPromptSingle("InsecureSkipVerify, true (Only recommended for Dev), false (Default)", false, "")
							if s == "false" {
								b.SslSkip = false
							} else {
								b.SslSkip = true
							}
							b.Connectivity(z)
						}
					}
				}
			}
		},
	}
}
