package amqp

func BuildRoutingKey(prefix string, initiator string) string {
	return prefix + "_" + initiator
}
