package connect

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func TestChannelConsume(t *testing.T) {
	cc := ChannelPoolConstructor{
		Connection: ConnectionConstructor{
			UserName: "rabbit",
			Password: "aneMicC7A9np",
			Address:  "localhost:5672",
		},
		MaxSize: 10,
	}.New()
	ch, err := cc.GetChannel()
	if err != nil {
		log.Panic(err)
	}

	msgs, err := ch.Consume(
		"hello", // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	var forever chan struct{}
	<-forever
}
func TestChannelPush(t *testing.T) {
	cc := ChannelPoolConstructor{
		Connection: ConnectionConstructor{
			UserName: "rabbit",
			Password: "aneMicC7A9np",
			Address:  "localhost:5672",
		},
		MaxSize: 10,
	}.New()
	ch, err := cc.GetChannel()
	if err != nil {
		log.Panic(err)
	}

	//q, err := Queue{
	//	Name:       "hello",
	//	Durable:    false,
	//	AutoDelete: false,
	//	Exclusive:  false,
	//	NoWait:     false,
	//	Args:       nil,
	//}.Declare(ch)
	//if err != nil {
	//	log.Panic(err)
	//}
	body := "Hello World!"
	err = ch.PublishWithContext(context.Background(),
		"",      // exchange
		"hello", // routing key
		false,   // mandatory
		false,   // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Panic(err)
	}
	log.Printf(" [x] Sent %s\n", body)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func TestSend(t *testing.T) {
	conn, err := amqp091.Dial("amqp091://rabbit:aneMicC7A9np@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

func TestConsume(t *testing.T) {
	conn, err := amqp091.Dial("amqp091://rabbit:aneMicC7A9np@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
