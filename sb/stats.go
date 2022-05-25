package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

type statsArgs struct {
	queueOrTopic string
	// or
	topic        string
	subscription string

	auth *Auth
	cmd  *cobra.Command
}

func newStatsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats [queue|topic subscription]",
		Short: "Show message statistics for a queue, topic or subscription.",
	}

	cmd.Args = cobra.RangeArgs(1, 2)

	statsArgs := &statsArgs{
		auth: AddAuth(cmd.Flags()),
		cmd:  cmd,
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			statsArgs.queueOrTopic = args[0]
		} else if len(args) == 2 {
			statsArgs.topic, statsArgs.subscription = args[0], args[1]
		}

		return statsCommand(statsArgs)
	}

	return cmd
}

func toInterface[T any](t *T) interface{} {
	if t == nil {
		return nil
	}

	return t
}

func statsCommand(args *statsArgs) error {
	if args.queueOrTopic == "" && args.topic == "" {
		args.cmd.Flags().PrintDefaults()
		return fmt.Errorf("you need to specify a queue, topic or subscription")
	}

	client, err := args.auth.NewAdminClient()

	if err != nil {
		return fmt.Errorf("failed to create a Service Bus admin client: %w", err)
	}

	if args.queueOrTopic != "" {
		type resp struct {
			err   error
			stats interface{}
		}

		ch := make(chan resp, 2)

		// a slight hack - we don't know if this is a queue or a topic
		// we'll try both - only one should work (or both fail because of some other issue)
		go func() {
			rt, err := client.GetTopicRuntimeProperties(context.Background(), args.queueOrTopic, nil)

			ch <- resp{
				err:   err,
				stats: toInterface(rt),
			}
		}()

		go func() {
			rt, err := client.GetQueueRuntimeProperties(context.Background(), args.queueOrTopic, nil)

			ch <- resp{
				err:   err,
				stats: toInterface(rt),
			}
		}()

		var msg string

		for i := 0; i < 2; i++ {
			resp := <-ch

			if resp.stats != nil {
				bytes, _ := json.MarshalIndent(resp.stats, "  ", "  ")
				msg = string(bytes)
				break
			}

			if resp.err == nil {
				msg = fmt.Sprintf("Failed to get queue/topic stats: not found")
			} else {
				msg = fmt.Sprintf("Failed to get queue/topic stats: %s", resp.err.Error())
			}
		}

		fmt.Printf("%s\n", msg)

	} else {
		rt, err := client.GetSubscriptionRuntimeProperties(context.Background(), args.topic, args.subscription, nil)

		if err == nil {
			if rt == nil {
				fmt.Printf("Failed to get subscription stats: not found\n")
			} else {
				bytes, _ := json.MarshalIndent(rt, "  ", "  ")
				fmt.Printf("%s\n", string(bytes))
			}
		}
	}

	return nil
}
