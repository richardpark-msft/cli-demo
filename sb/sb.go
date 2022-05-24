package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "sb",
	}

	rootCmd.AddCommand(newSendCommand())
	rootCmd.AddCommand(newReceiveCommand())
	rootCmd.AddCommand(newPeekCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		os.Exit(1)
	}
}
