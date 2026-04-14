package main

import (
	"context"
)

type PollingChatClient interface {
	PollEvents() []ChatEvent
}

type ChatClient interface {
	PollingChatClient
}

func NewChatClient(ctx context.Context) ChatClient {
	return nil
}

// --- mocks ---

type mockedChatEvent struct {
	message string
}

type mockChatClient struct {
	defaultMessage string
}

func (ce *mockedChatEvent) Message() string {
	return ce.message
}

func (c *mockChatClient) PollEvents() []ChatEvent {
	return []ChatEvent{
		mockChatEvent(c.defaultMessage),
		mockChatEvent(c.defaultMessage),
		mockChatEvent(c.defaultMessage),
	}
}

func mockChatEvent(message string) ChatEvent {
	return &mockedChatEvent{
		message: message,
	}
}

func MockChatClient(ctx context.Context) ChatClient {
	return &mockChatClient{
		defaultMessage: "Mock",
	}
}
