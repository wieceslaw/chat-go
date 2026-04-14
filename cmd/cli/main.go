package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Fprintln(os.Stdout, "Hello")

	ctx := context.Background()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	client := MockChatClient(ctx)

	channel := pollChatEvents(ctx, client)
	go func() {
		for event := range channel {
			println(event.Message())
		}
	}()

	select {
	case <-ctx.Done():
	case <-time.After(time.Second * 10):
	}
	stop()
}
