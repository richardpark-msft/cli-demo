package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("Usage: sb (send|receive|manage)\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "send":
		if err := sendCommand(os.Args[2:]); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			os.Exit(1)
		}
	case "receive":
		break
	}
}

func sendCommand(args []string) error {
	fs := flag.NewFlagSet("send", flag.ExitOnError)

	envFile := fs.String("env", ".env", ".env file to load before running. Empty string disables .env file loading.")
	entityPath := fs.String("entity", "", "Entity to send messages to (topic or queue)")
	auth := AddAuth(fs)

	isBody := fs.Bool("body", true, "Causes the entire contents of stdin to be used as the message Body.")

	_ = fs.Parse(os.Args[2:])

	if *entityPath == "" {
		return fmt.Errorf("no queue/topic specified with -entity")
	}

	if *envFile != "" {
		if err := godotenv.Load(*envFile); err != nil {
			return fmt.Errorf("failed to load .env file '%s': %s", *envFile, err.Error())
		}
	}

	client, err := auth.NewClient()

	if err != nil {
		return err
	}

	sender, err := client.NewSender(*entityPath)

	if err != nil {
		return err
	}

	// read the message in from stdin
	if *isBody {
		bytes, err := ioutil.ReadAll(os.Stdin)

		if err != nil {
			return err
		}

		// send the entire message as a body
		msg := &azservicebus.Message{
			Body: bytes,
		}

		return sender.SendMessage(context.Background(), msg)
	} else {
		return fmt.Errorf("sending full messages is not yet implemented")
	}
}
