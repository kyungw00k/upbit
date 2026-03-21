package websocket

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Format 구독 메시지 포맷 타입
type Format string

const (
	// FormatDefault 기본 포맷 (모든 필드 포함)
	FormatDefault Format = "DEFAULT"
	// FormatSimple 간소화 포맷
	FormatSimple Format = "SIMPLE"
)

// SubscriptionType 구독 데이터 타입
type SubscriptionType struct {
	Type           string   `json:"type"`
	Codes          []string `json:"codes,omitempty"`
	IsOnlySnapshot bool     `json:"is_only_snapshot,omitempty"`
	IsOnlyRealtime bool     `json:"is_only_realtime,omitempty"`
}

// BuildSubscribeMessage 구독 요청 메시지 생성
// ticket은 자동으로 UUID가 생성됨, format은 DEFAULT
func BuildSubscribeMessage(types []SubscriptionType) ([]byte, error) {
	ticket := uuid.New().String()
	return BuildSubscribeMessageWithOptions(ticket, types, FormatDefault)
}

// BuildSubscribeMessageWithTicket 지정된 ticket으로 구독 요청 메시지 생성
func BuildSubscribeMessageWithTicket(ticket string, types []SubscriptionType) ([]byte, error) {
	return BuildSubscribeMessageWithOptions(ticket, types, FormatDefault)
}

// BuildSubscribeMessageWithOptions ticket, format 지정으로 구독 요청 메시지 생성
// W-1: format 필드를 포함하여 Upbit API 문서에 맞는 구독 메시지 생성
func BuildSubscribeMessageWithOptions(ticket string, types []SubscriptionType, format Format) ([]byte, error) {
	// 배열 구성: [{"ticket": "uuid"}, {"type": "ticker", "codes": [...]}, ..., {"format": "DEFAULT"}]
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
