/*
Copyright © 2024 moti <motohiko.ave@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:  "nostr-crawler",
	Long: "crawls nostr notes and publish them to a relay",
	RunE: func(cmd *cobra.Command, args []string) error {
		froms, err := cmd.Flags().GetStringArray("from")
		if err != nil {
			return err
		}
		if len(froms) == 0 {
			return fmt.Errorf("no relay 'from' specified")
		}
		if len(froms) > 1 {
			return fmt.Errorf("currently only one relay can be specified in 'from'")
		}

		tos, err := cmd.Flags().GetStringArray("to")
		if err != nil {
			return err
		}
		if len(tos) == 0 {
			return fmt.Errorf("no relay 'to' specified")
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-signals
			cancel()
		}()

		sub := getSubscription(ctx, froms[0])

		for {
			select {
			case ev := <-sub.Events:
				publish(tos, *ev)
			case <-ctx.Done():
				return nil
			}
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func getSubscription(ctx context.Context, from string) *nostr.Subscription {
	relay, err := nostr.RelayConnect(ctx, from)
	if err != nil {
		panic(err)
	}

	var filters nostr.Filters
	since := nostr.Timestamp(time.Now().Unix())
	filters = []nostr.Filter{{
		Since: &since,
	}}

	sub, err := relay.Subscribe(ctx, filters)
	if err != nil {
		fmt.Println(err)
	}
	return sub
}

func publish(tos []string, ev nostr.Event) {
	type relayurl string
	fmt.Printf("%+v\n", ev)
	for _, to := range tos {
		ctx := context.WithValue(context.Background(), relayurl("relay"), to)
		relay, e := nostr.RelayConnect(ctx, to)
		if e != nil {
			fmt.Println(e)
			continue
		}
		fmt.Println("posting to: ", to)
		relay.Publish(ctx, ev)
	}
}

func init() {
	rootCmd.Flags().StringArrayP("from", "f", []string{}, "relay to crawl")
	rootCmd.Flags().StringArrayP("to", "t", []string{}, "relay to publish")
}
