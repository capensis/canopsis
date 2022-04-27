package main

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	liblog "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	daemonName = "amqp2tty"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	amqpConnection, err := amqp.NewConnection(liblog.NewLogger(true), 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := amqpConnection.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	ch, err := amqpConnection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := ch.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	exchange := canopsis.CanopsisEventsExchange

	queue, err := ch.QueueDeclare(
		daemonName,
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(
		queue.Name,
		"#",
		exchange,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		queue.Name,
		daemonName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s started\n", daemonName)
	defer fmt.Printf("\n%s closed\n", daemonName)

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-msgs:
			if !ok {
				log.Fatal("the rabbitmq channel has been closed")
			}

			fmt.Printf("%s %s New message:\n%s@%s:\n\t%s\n",
				time.Now().Format("2006-01-02T15:04:05.999999999Z07:00"),
				daemonName, d.RoutingKey, d.Exchange, d.Body)
		}
	}
}
