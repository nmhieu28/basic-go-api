package cqrs

import (
	"errors"
	"fmt"
)

type Mediator interface {
	Register(handler interface{}, commands ...interface{})
	Send(command interface{}) (interface{}, error)
}
type mediator struct {
	handlers map[string]interface{}
}

func NewMediator() *mediator {
	return &mediator{handlers: make(map[string]interface{})}
}
func (m *mediator) Send(command Command) (interface{}, error) {
	commandName := typeOf(command)
	handler, ok := m.handlers[commandName]
	if !ok {
		return nil, errors.New("handler not found")
	}

	switch h := handler.(type) {
	case func(interface{}) (interface{}, error):
		return h(command)
	default:
		return nil, errors.New("invalid handler type")
	}
}
func (m *mediator) Register(handler interface{}, commands ...any) error {
	for _, command := range commands {
		typeName := typeOf(command)
		if _, ok := m.handlers[typeName]; ok {
			return fmt.Errorf("Duplicate command handler registration with command bus for command of type: %s", typeName)
		}
		m.handlers[typeName] = handler
	}
	return nil
}
