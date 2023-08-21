package adepter

import (
	"context"
	"log"

	"github.com/meowalien/RabbitGather-proto/go/interest"
	"github.com/meowalien/RabbitGather-proto/go/share"
	"github.com/meowalien/go-meowalien-lib/errs"
	"github.com/rabbitmq/amqp091-go"

	"github.com/meowalien/RabbitGather-interest-crawler.git/fremwork/connect"
	"github.com/meowalien/RabbitGather-interest-crawler.git/lib"
)

type MQConsumer struct {
	Queue   connect.Queue
	Consume connect.Consume
	Handler MQHandler
}

func (a *adepter) RabbitMQ(cp connect.ChannelPool) {
	ch, err := cp.GetChannel()
	if err != nil {
		panic(err)
	}
	consumers := []MQConsumer{
		{
			Queue: connect.Queue{
				Name:       interest.InterestCrawler_Crawl_FullMethodName,
				Durable:    false,
				AutoDelete: false,
				Exclusive:  false,
				NoWait:     false,
				Args:       nil,
			},
			Consume: connect.Consume{
				Consumer:  "",
				AutoAck:   false,
				Exclusive: false,
				NoLocal:   false,
				NoWait:    false,
				Args:      nil,
			},
			Handler: a.rabbitMQCrewHandler,
		},
	}
	for _, consumer := range consumers {
		var q amqp091.Queue
		q, err = consumer.Queue.Declare(ch)
		if err != nil {
			panic(err)
		}
		var deliveryCh <-chan amqp091.Delivery
		deliveryCh, err = consumer.Consume.Consume(ch, q.Name)
		if err != nil {
			panic(err)
		}
		for delivery := range deliveryCh {
			a.autoRecoverOrAck(delivery, consumer.Handler)
		}
	}
}
func (a *adepter) rabbitMQCrewHandler(ctx context.Context, ch *amqp091.Channel, delivery amqp091.Delivery) (err error) {
	req, _, err := decodeMQMessage[interest.CrawlRequest](delivery)
	_, err = a.md.Crawl(ctx, req)
	if err != nil {
		return
	}
	return
}

func decodeMQMessage[T any](delivery amqp091.Delivery) (res *T, encoding int32, err error) {
	var ok1 bool
	encoding, ok1 = delivery.Headers[EncodingKey].(int32)
	if !ok1 {
		err = errs.New("FullMethodName not found")
		return
	}

	err = lib.DecodeMessage(share.Encoding(encoding), delivery.Body, res)
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}

const EncodingKey = "Encoding"

func (a *adepter) autoRecoverOrAck(delivery amqp091.Delivery, fc MQHandler) {
	defer func() {
		if e := recover(); e != nil {
			e1 := delivery.Reject(false)
			if e1 != nil {
				e1 = errs.New(e1)
				log.Println(e1)
				return
			}
			e = errs.New(e)
			log.Println(e)
		}
	}()

	e := fc(context.Background(), nil, delivery)

	if e != nil {
		e1 := delivery.Reject(true)
		if e1 != nil {
			e1 = errs.New(e1)
			log.Println(e1)
			return
		}
	}
	e2 := delivery.Ack(false)
	if e2 != nil {
		e2 = errs.New(e2)
		log.Println(e2)
		return
	}
}

type MQHandler func(ctx context.Context, ch *amqp091.Channel, delivery amqp091.Delivery) error
