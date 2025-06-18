package main

import (
	"fmt"
	"s2a/cmd/s2a"
)

func main() {
	rootCmd := s2a.NewRootCmd()
	rootExecErr := rootCmd.Execute()
	if rootExecErr != nil {
		fmt.Println(rootExecErr)
		return
	}
}
