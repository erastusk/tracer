/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github/erastusk/tracer/cmd"
	"github/erastusk/tracer/internal/utils"
)

func main() {
	utils.Done = make(chan struct{})
	cmd.Execute()
}
