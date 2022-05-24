package internal

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type FormatMessageArgs struct {
	Message      *azservicebus.ReceivedMessage
	OneLine      bool
	AssumeString bool
}

func FormatMessage(args FormatMessageArgs) ([]byte, error) {
	formattableMessage := map[string]interface{}{}

	if args.AssumeString {
		formattableMessage["Body"] = string(args.Message.Body)
	} else {
		formattableMessage["Body"] = args.Message.Body
	}

	// copy over the common fields
	if args.Message.ApplicationProperties != nil {
		formattableMessage["ApplicationProperties"] = args.Message.ApplicationProperties
	}

	formattableMessage["SequenceNumber"] = args.Message.SequenceNumber
	formattableMessage["DeliveryCount"] = args.Message.DeliveryCount
	formattableMessage["EnqueuedTime"] = args.Message.EnqueuedTime
	formattableMessage["ExpiresAt"] = args.Message.ExpiresAt

	if args.Message.Subject != nil {
		formattableMessage["Subject"] = args.Message.Subject
	}

	if args.Message.SessionID != nil {
		formattableMessage["SessionID"] = args.Message.SessionID
	}

	formattableMessage["MessageID"] = args.Message.MessageID

	var bytes []byte
	var err error

	if args.OneLine {
		bytes, err = json.Marshal(formattableMessage)
	} else {
		bytes, err = json.MarshalIndent(formattableMessage, "  ", "  ")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to JSON.Marshal message with MessageID %s: %w", args.Message.MessageID, err)
	}

	return bytes, nil
}
