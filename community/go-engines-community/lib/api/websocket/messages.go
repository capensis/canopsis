package websocket

const (
	ClientMessageClientPing = iota
	ClientMessageJoin
	ClientMessageLeave
	ClientMessageAuth
	ClientMessageInfo
)

const (
	ServerMessageClientPong = iota
	ServerMessageInfo
	ServerMessageError
	ServerMessageCloseRoom
	ServerMessageAuthSuccess
)

type ClientMessage struct {
	Type    int    `json:"type"`
	Room    string `json:"room"`
	Token   string `json:"token"`
	Payload any    `json:"payload"`
}

type ServerMessage struct {
	Type    int    `json:"type"`
	Room    string `json:"room,omitempty"`
	Payload any    `json:"payload,omitempty"`
	Error   int    `json:"error,omitempty"`
}

type ConnectionInfo struct {
	ID     string
	UserID string
}

type ConnectionAuthInfo struct {
	ID     string
	UserID string
	Token  string
}
