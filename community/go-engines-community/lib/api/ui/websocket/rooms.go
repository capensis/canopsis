package websocket

const (
	RoomLoggedUserCount   = "logged-user-count"
	RoomBroadcastMessages = "broadcast-messages"

	RoomAlarmsGroup       = "alarms"
	RoomAlarmDetailsGroup = "alarm-details"

	RoomHealthcheck       = "healthcheck"
	RoomHealthcheckStatus = "healthcheck-status"
	RoomMessageRates      = "message-rates"
	RoomIcons             = "icons"
	RoomNotifications     = "notifications"
	RoomPbhPatterns       = "pbehavior-patterns"
)

func GroupRoom(g, id string) string {
	return g + groupDelimiter + id
}
