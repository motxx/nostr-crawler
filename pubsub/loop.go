package pubsub

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nbd-wtf/go-nostr"
)

func Loop(froms []string, tos []string) error {
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
		relay.Publish(ctx, ev)
	}
}
