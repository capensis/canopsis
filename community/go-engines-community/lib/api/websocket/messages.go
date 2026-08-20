package websocket

//go:generate go tool github.com/mailru/easyjson/easyjson -no_std_marshalers

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
	ServerMessageJoined
	ServerMessageLeft
)

// ClientMessage
//
// easyjson:json
type ClientMessage struct {
	Type    int    `json:"type"`
	Room    string `json:"room"`
	Token   string `json:"token"`
	Payload any    `json:"payload"`
}

// ServerMessage
//
// easyjson:json
type ServerMessage struct {
	Type    int    `json:"type"`
	Room    string `json:"room,omitempty"`
	Payload any    `json:"payload,omitempty"`
	Error   int    `json:"error,omitempty"`
}

type User struct {
	ID     string
	Locale string
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
