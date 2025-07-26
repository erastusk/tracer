package kafka

import (
	"github/erastusk/tracer/cmd/kafka/confluent"
	"github/erastusk/tracer/cmd/kafka/msk"

	"github.com/spf13/cobra"
)

// redisCmd represents the kafka command
var KafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "kafka connectivity check",
	Long:  "kafka connectivity check",
}

func init() {
	KafkaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	KafkaCmd.AddCommand(confluent.ConfluentCmd)
	KafkaCmd.AddCommand(msk.MskCmd)
}
