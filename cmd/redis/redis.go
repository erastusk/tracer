/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package redis

import "github/erastusk/tracer/internal/Cmd"

// redisCmd represents the redis command
var (
	a        = Cmd.NewRedisType("redis")
	RedisCmd = a.Cmd()
)

func init() {}
