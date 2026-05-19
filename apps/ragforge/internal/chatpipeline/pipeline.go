package chatpipeline

import (
	"context"
)

type EventManager struct {
	listeners map[EventType][]Plugin
	handlers  map[EventType]func(context.Context, EventType, *ChatManage) *PluginError
}

func NewEventManager() *EventManager {
	return &EventManager{
		listeners: make(map[EventType][]Plugin),
		handlers:  make(map[EventType]func(context.Context, EventType, *ChatManage) *PluginError),
	}
}

func (e *EventManager) Register(plugin Plugin) {
	if e.listeners == nil {
		e.listeners = make(map[EventType][]Plugin)
	}
	if e.handlers == nil {
		e.handlers = make(map[EventType]func(context.Context, EventType, *ChatManage) *PluginError)
	}
	for _, eventType := range plugin.ActivationEvents() {
		e.listeners[eventType] = append(e.listeners[eventType], plugin)
		e.handlers[eventType] = e.buildHandler(e.listeners[eventType])
	}
}

func (e *EventManager) buildHandler(plugins []Plugin) func(
	ctx context.Context, eventType EventType, chatManage *ChatManage,
) *PluginError {
	next := func(context.Context, EventType, *ChatManage) *PluginError { return nil }
	for i := len(plugins) - 1; i >= 0; i-- {
		current := plugins[i]
		prevNext := next
		next = func(ctx context.Context, eventType EventType, chatManage *ChatManage) *PluginError {
			return current.OnEvent(ctx, eventType, chatManage, func() *PluginError {
				return prevNext(ctx, eventType, chatManage)
			})
		}
	}
	return next
}

func (e *EventManager) Trigger(ctx context.Context,
	eventType EventType, chatManage *ChatManage,
) *PluginError {
	if handler, ok := e.handlers[eventType]; ok {
		return handler(ctx, eventType, chatManage)
	}
	return nil
}
