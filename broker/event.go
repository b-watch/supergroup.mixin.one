package broker

import (
	"encoding/json"
)

type Event struct {
	Topic string                 `json:"topic,omitempty"`
	Body  json.RawMessage        `json:"body,omitempty"`
	Meta  map[string]interface{} `json:"meta,omitempty"`
}

func (e *Event) WithMeta(key string, value interface{}) {
	if e.Meta == nil {
		e.Meta = make(map[string]interface{})
	}

	e.Meta[key] = value
}

func (e *Event) GetMeta(key string) (interface{}, bool) {
	if e.Meta == nil {
		return nil, false
	}

	v, ok := e.Meta[key]
	return v, ok
}
