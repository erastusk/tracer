package Cmd

import (
	"log"

	_cache "github/erastusk/tracer/internal/cache"
	"github/erastusk/tracer/internal/confluent"
	"github/erastusk/tracer/internal/msk"
	"github/erastusk/tracer/internal/prompts"
	"github/erastusk/tracer/internal/types"
	"github/erastusk/tracer/internal/utils"

	"github.com/spf13/cobra"
)

type CobraCMD interface {
	Cmd() *cobra.Command
}

type KafkaType struct {
	conf  types.KafkaApps
	kType string
	cache *_cache.Cache
}

func NewKafkaType(s string) *KafkaType {
	var c types.KafkaApps
	var newcache *_cache.Cache
	switch s {
	case "confluent":
		conf, err := confluent.NewConfluentOauth()
		if err != nil {
			log.Fatal(err)
		}
		c = conf
	case "msk":
		newcache = _cache.NewCache(s)
		c = msk.NewMsk()
	}
	return &KafkaType{
		conf:  c,
		kType: s,
		cache: newcache,
	}
}

func (k *KafkaType) Cmd() *cobra.Command {
	count := 0
	var c types.PromptOptions
	utils.InitColors()
	b := utils.NewKafkaApp(k.conf)
	tcp := utils.NewTestConnectivityImpl()
	return &cobra.Command{
		Use:   k.kType,
		Short: k.kType + " kafka Connectivity diagnosis tool",
		Run: func(cmd *cobra.Command, args []string) {
			options := []string{"TCP Connnectivity test", "Broker connectivity test", "List topics", "Produce", "Consume"}
			utils.Green.Printf("Processing commands for %s kafka\n", k.kType)
			for {
				if k.kType == "msk" {
					count++
					if count > 0 {
						c = k.cache.LoadCache()
					}
				}
				select {
				case <-utils.Done:
					utils.Green.Println("exiting....")
					return
				default:
					switch prompts.GetUserPrompt(options) {
					case "TCP Connnectivity test":
						err := tcp.TCPDial(prompts.GetUserPromptSingle("Please enter broker", false, c.Endpoint))
						if err != nil {
							utils.Red.Printf("TCP connection failed: %v\n", err)
						}
						utils.Green.Println("TCP connection succeeded to -> ", c.Endpoint)
					case "Broker connectivity test":
						if k.kType == "msk" {
							z, err := prompts.GetPrompts(c, utils.Options)
							if err != nil {
								utils.Red.Println(err)
							}
							k.cache.SaveCache(z)
							err = b.App.Connectivity(z)
							if err != nil {
								utils.Red.Println("Broker connectivity was unsuccessful...")
							}
							break
						}
						err := b.App.Connectivity(types.PromptOptions{})
						if err != nil {
							utils.Red.Println("Broker connectivity was unsuccessful...")
						}
						utils.Green.Println("Broker connectivity was successful!")
					case "List topics":
						if k.kType == "msk" {
							z, err := prompts.GetPrompts(c, utils.Options)
							if err != nil {
								utils.Red.Println(err)
							}
							k.cache.SaveCache(z)
							err = b.App.ListTopics(z)
							if err != nil {
								utils.Red.Println("List Topics was unsuccessful...")
							}
							break
						}
						err := b.App.ListTopics(types.PromptOptions{})
						if err != nil {
							utils.Red.Println("List Topics was unsuccessful...")
						}
					case "Produce":
						var z types.PromptOptions
						if k.kType == "msk" {
							zp, err := prompts.GetPrompts(c, utils.Options)
							if err != nil {
								utils.Red.Println(err)
								break
							}
							z = zp
							k.cache.SaveCache(z)

						}
						t := prompts.GetUserPromptSingle("Please enter topic", false, "")
						if t == "" {
							utils.Red.Println("Topic required")
						} else {
							err := b.App.Produce(z, t)
							if err != nil {
								utils.Red.Println("Producer unsuccessful...")
							} else {
								utils.Green.Printf("Successfully produced to topic %s\n", t)
							}
						}
					case "Consume":
						var z types.PromptOptions
						if k.kType == "msk" {
							zp, err := prompts.GetPrompts(c, utils.Options)
							if err != nil {
								utils.Red.Println(err)
								break
							}
							z = zp
							k.cache.SaveCache(z)

						}
						t := prompts.GetUserPromptSingle("Please enter topic", false, "")
						if t == "" {
							utils.Red.Println("Topic required")
						} else {
							err := b.App.Consume(z, t)
							if err != nil {
								utils.Red.Println("Consumer unsuccessful...", err)
							} else {
								utils.Green.Printf("Successfully consumed from topic %s\n", t)
							}
						}
					}
				}
			}
		},
	}
}
