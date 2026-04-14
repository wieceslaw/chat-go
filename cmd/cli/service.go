package main

import (
	"context"
	"time"
)

type ChatEvent interface {
	Message() string
}

type ChatService interface {
	GetChatEvents() <-chan ChatEvent
}

// polling integration implementation
func pollChatEvents(ctx context.Context, client PollingChatClient) <-chan ChatEvent {
	channel := make(chan ChatEvent)
	delay := time.Second * 1

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(channel)
				return
			case <-time.After(delay):
				for _, event := range client.PollEvents() {
					channel <- event
				}
			}
		}
	}()

	return channel
}
