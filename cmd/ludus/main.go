package main

import (
	"fmt"
	"github.com/tanema/ludus/cmd"
)

func main() {
	if err := cmd.LudusCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
