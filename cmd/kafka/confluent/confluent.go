/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package confluent

import "github/erastusk/tracer/internal/Cmd"

// confluentCmd represents the confluent command
var (
	a            = Cmd.NewKafkaType("confluent")
	ConfluentCmd = a.Cmd()
)

func init() {}
