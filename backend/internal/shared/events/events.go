package events

// TODO: Implement event system
// This will handle:
// - Domain events publishing
// - Event subscribers
// - Event sourcing (optional)
// - Async event processing

type Event interface {
	EventType() string
	EventData() interface{}
	EventTime() string
}

type EventBus interface {
	Publish(event Event) error
	Subscribe(eventType string, handler EventHandler) error
	Unsubscribe(eventType string, handler EventHandler) error
}

type EventHandler interface {
	Handle(event Event) error
}

// TODO: Implement common events
// - TenantCreated
// - TenantUpdated
// - UserRegistered
// - UserLoggedIn
// - ProductCreated
// - ProductUpdated
// - OrderPlaced
// - OrderUpdated
// - PaymentProcessed
// - PaymentFailed
// - NotificationSent
