/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package msk

import "github/erastusk/tracer/internal/Cmd"

// mskCmd represents the msk command
var (
	a      = Cmd.NewKafkaType("msk")
	MskCmd = a.Cmd()
)

func init() {
}
