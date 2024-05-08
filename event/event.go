package event

import (
	"errors"
	"fmt"
)

type Event interface {
	Name() string
	Deserialize(data []byte) error
}

var eventMap = map[string]func() Event{
	"fruit.orange": func() Event { return &OrangeEvent{} },
}

type Factory struct {
}

func (f Factory) Create(name string) (Event, error) {
	if event, ok := eventMap[name]; ok {
		return event(), nil
	}

	return nil, errors.New(fmt.Sprintf("Event not found: %s", name))
}

func NewFactory() *Factory {
	return &Factory{}
}
