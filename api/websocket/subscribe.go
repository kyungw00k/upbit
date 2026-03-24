package websocket

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Format is the subscription message format type.
type Format string

const (
	// FormatDefault is the default format (all fields included).
	FormatDefault Format = "DEFAULT"
	// FormatSimple is the simplified format.
	FormatSimple Format = "SIMPLE"
)

// SubscriptionType specifies the data type to subscribe to.
type SubscriptionType struct {
	Type           string   `json:"type"`
	Codes          []string `json:"codes,omitempty"`
	IsOnlySnapshot bool     `json:"is_only_snapshot,omitempty"`
	IsOnlyRealtime bool     `json:"is_only_realtime,omitempty"`
}

// BuildSubscribeMessage creates a subscription request message.
// A UUID ticket is generated automatically; format defaults to DEFAULT.
func BuildSubscribeMessage(types []SubscriptionType) ([]byte, error) {
	ticket := uuid.New().String()
	return BuildSubscribeMessageWithOptions(ticket, types, FormatDefault)
}

// BuildSubscribeMessageWithTicket creates a subscription request message with the specified ticket.
func BuildSubscribeMessageWithTicket(ticket string, types []SubscriptionType) ([]byte, error) {
	return BuildSubscribeMessageWithOptions(ticket, types, FormatDefault)
}

// BuildSubscribeMessageWithOptions creates a subscription request message with the specified ticket and format.
// W-1: Includes the format field to produce a subscription message conforming to the Upbit API spec.
func BuildSubscribeMessageWithOptions(ticket string, types []SubscriptionType, format Format) ([]byte, error) {
	// Array layout: [{"ticket": "uuid"}, {"type": "ticker", "codes": [...]}, ..., {"format": "DEFAULT"}]
	msg := make([]interface{}, 0, len(types)+2)

	// Ticket Object
	msg = append(msg, map[string]string{"ticket": ticket})

	// Data Type Objects
	for _, t := range types {
		obj := map[string]interface{}{
			"type": t.Type,
		}
		if len(t.Codes) > 0 {
			obj["codes"] = t.Codes
		}
		if t.IsOnlySnapshot {
			obj["is_only_snapshot"] = true
		}
		if t.IsOnlyRealtime {
			obj["is_only_realtime"] = true
		}
		msg = append(msg, obj)
	}

	// Format Object
	if format != "" {
		msg = append(msg, map[string]string{"format": string(format)})
	}

	return json.Marshal(msg)
}
