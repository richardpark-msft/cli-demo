package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/parkplusplus/cli/sb/internal"
	"github.com/spf13/cobra"
)

type peekArgs struct {
	oneLine        bool
	timeout        time.Duration
	count          int
	sequenceNumber int
	auth           *Auth

	queue string
	// or
	topic        string
	subscription string

	cmd *cobra.Command
}

func newPeekCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "peek [queue|topic subscription] [flags]",
		Short: "Peek messages from a queue or a topic and subscription.",
	}

	peekArgs := &peekArgs{
		auth: AddAuth(cmd.Flags()),
		cmd:  cmd,
	}

	cmd.Args = cobra.RangeArgs(1, 2)
	cmd.Flags().BoolVar(&peekArgs.oneLine, "oneline", true, "Print each message as a single line.")
	cmd.Flags().IntVarP(&peekArgs.count, "count", "c", 1, "Max messages to peek.")
	cmd.Flags().IntVarP(&peekArgs.sequenceNumber, "sequence-number", "s", 0, "Starting sequence number")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			peekArgs.queue = args[0]
		} else if len(args) == 2 {
			peekArgs.topic, peekArgs.subscription = args[0], args[1]
		}

		return peekCommand(peekArgs)
	}

	return cmd
}

func peekCommand(args *peekArgs) error {
	client, err := args.auth.NewClient()

	if err != nil {
		return fmt.Errorf("failed to create a Service Bus client: %w", err)
	}

	defer client.Close(context.Background())

	var receiver *azservicebus.Receiver

	// when peeking messages you don't actually utilize the receive mode since peeked
	// messages are not settled.
	if args.queue != "" {
		receiver, err = client.NewReceiverForQueue(args.queue, nil)
	} else {
		receiver, err = client.NewReceiverForSubscription(args.topic, args.subscription, nil)
	}

	if err != nil {
		return err
	}

	var options *azservicebus.PeekMessagesOptions

	if args.sequenceNumber != 0 {
		options = &azservicebus.PeekMessagesOptions{
			FromSequenceNumber: to.Ptr(int64(args.sequenceNumber)),
		}
	}

	peekedMessages, err := receiver.PeekMessages(context.Background(), args.count, options)

	if err != nil {
		return err
	}

	for _, m := range peekedMessages {
		bytes, err := internal.FormatMessage(internal.FormatMessageArgs{
			Message:      m,
			OneLine:      args.oneLine,
			AssumeString: true,
		})

		if err != nil {
			return fmt.Errorf("failed to format message: %s", err)
		}

		fmt.Printf("%s\n", string(bytes))
	}

	return nil
}
