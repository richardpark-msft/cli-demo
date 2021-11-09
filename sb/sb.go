package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
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

	// authentication methods
	ns := fs.String("ns", "", "Namespace (assumes DefaultAzureCredential)")
	cs := fs.String("csenv", "SERVICEBUS_CONNECTION_STRING", "Environment variable that contains a connection string")
	queueOrTopic := fs.String("entity", "", "Entity to send messages to (topic or queue)")
	isBody := fs.Bool("body", true, "Causes the entire contents of stdin to be used as the message Body.")

	_ = fs.Parse(os.Args[2:])

	if *queueOrTopic == "" {
		return fmt.Errorf("no queue/topic specified with -entity")
	}

	_ = godotenv.Load()

	var client *azservicebus.Client

	if *ns != "" {
		dac, err := azidentity.NewDefaultAzureCredential(nil)

		if err != nil {
			return err
		}

		client, err = azservicebus.NewClient(*ns, dac, nil)

		if err != nil {
			return err
		}
	} else {
		cs := os.Getenv(*cs)

		if cs == "" {
			return fmt.Errorf("no connection string in environment variable %s", cs)
		}

		var err error
		client, err = azservicebus.NewClientFromConnectionString(cs, nil)

		if err != nil {
			return err
		}
	}

	sender, err := client.NewSender(*queueOrTopic)

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
