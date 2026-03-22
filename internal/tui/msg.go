package tui

// WsMsg WebSocket에서 수신한 원시 데이터 메시지
type WsMsg struct {
	Data []byte
}

// ErrMsg 에러 메시지
type ErrMsg struct {
	Err error
}

// ServerErrMsg 서버 에러 메시지
type ServerErrMsg struct {
	Message string
}
