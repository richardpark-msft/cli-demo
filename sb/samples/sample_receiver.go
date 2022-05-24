package main

import (
	"context"
	"log"
	"os"

	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Failed to load .env file: %s", err.Error())
	}

	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING")

	client, err := azservicebus.NewClientFromConnectionString(cs, nil)

	if err != nil {
		log.Fatalf("Failed to create client from connection string: %s", err.Error())
	}

	queueName := "demo"

	receiver, err := client.NewReceiverForQueue(queueName, &azservicebus.ReceiverOptions{
		ReceiveMode: azservicebus.ReceiveModeReceiveAndDelete,
	})

	if err != nil {
		log.Fatalf("Failed to create receiver: %s", err.Error())
	}

	log.Printf("Listening for messages...")

	for {
		messages, err := receiver.ReceiveMessages(context.Background(), 10, nil)

		if err != nil {
			log.Fatalf("Failed to receive messages: %s", err.Error())
		}

		for _, m := range messages {
			var data *struct {
				Planet string
			}

			if err := json.Unmarshal(m.Body, &data); err != nil {
				log.Printf("Failed to unmarshal data from %s: %s", string(m.Body), err)
			} else {
				switch data.Planet {
				case "Mars":
					log.Printf("Launching, destination Mars...")
				case "Earth":
					log.Printf("Driving, destination Earth")
				default:
					log.Printf("Unknown travel method, destination %s", data.Planet)
				}
			}
		}
	}
}
