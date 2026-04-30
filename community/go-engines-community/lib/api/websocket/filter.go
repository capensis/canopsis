package websocket

import "slices"

type Filter interface {
	room() string
	match(conn *connection) bool
}

func ToRoom(room string) Filter {
	return &toRoomFilter{r: room}
}

func ToUser(room string, userIDs ...string) Filter {
	return &toUserFilter{
		userIDs: userIDs,
		r:       room,
	}
}

func ToConnection(room string, connIDs ...string) Filter {
	return &toConnFilter{
		connIDs: connIDs,
		r:       room,
	}
}

type toRoomFilter struct {
	r string
}

func (f *toRoomFilter) room() string {
	return f.r
}

func (f *toRoomFilter) match(conn *connection) bool {
	return conn.inRoom(f.r)
}

type toUserFilter struct {
	r       string
	userIDs []string
}

func (f *toUserFilter) room() string {
	return f.r
}

func (f *toUserFilter) match(conn *connection) bool {
	user, _ := conn.getAuth()

	return slices.Contains(f.userIDs, user.ID) && conn.inRoom(f.r)
}

type toConnFilter struct {
	r       string
	connIDs []string
}

func (f *toConnFilter) room() string {
	return f.r
}

func (f *toConnFilter) match(conn *connection) bool {
	return slices.Contains(f.connIDs, conn.id) && conn.inRoom(f.r)
}
