package connect

import (
	"fmt"
	"log"
	"sync"

	"github.com/meowalien/go-meowalien-lib/errs"
	"github.com/rabbitmq/amqp091-go"
)

type ConnectionConstructor struct {
	UserName string
	Password string
	Address  string
}

func (a *ConnectionConstructor) Connect() (cn *amqp091.Connection, err error) {
	return amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s/", a.UserName, a.Password, a.Address))
}

type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp091.Table
}

func (q Queue) Declare(ch *amqp091.Channel) (amqp091.Queue, error) {
	return ch.QueueDeclare(
		q.Name,
		q.Durable,
		q.AutoDelete,
		q.Exclusive,
		q.NoWait,
		q.Args,
	)
}

type ChannelPool interface {
	GetChannel() (ch *amqp091.Channel, err error)
}

type ChannelPoolConstructor struct {
	Connection ConnectionConstructor
	MaxSize    int
}

func (c ChannelPoolConstructor) New() ChannelPool {
	return &channelPool{ChannelPoolConstructor: c, channelPool: make(chan *amqp091.Channel, c.MaxSize)}
}

type channelPool struct {
	ChannelPoolConstructor
	conn        *amqp091.Connection
	connLocker  sync.RWMutex
	connCond    *sync.Cond
	channelPool chan *amqp091.Channel
}

func (c *channelPool) GetChannel() (ch *amqp091.Channel, err error) {
	select {
	case ch = <-c.channelPool:
		return
	default:
	}
	ch, err = c.getChannel()
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}
func (c *channelPool) getChannel() (ch *amqp091.Channel, err error) {
	conn, err := c.getConnection()
	if err != nil {
		err = errs.New(err)
		return
	}

	ch, err = conn.Channel()
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}
func (c *channelPool) Put(ch *amqp091.Channel) {
	select {
	case c.channelPool <- ch:
	default:
		// channel pool is full, close channel
		err := ch.Close()
		if err != nil {
			log.Panicln("error close channel", err)
		}
	}
}

func (c *channelPool) getConnection() (conn *amqp091.Connection, err error) {
	c.connLocker.RLock()
	if c.conn != nil {
		// connection is not nil, return it
		defer c.connLocker.RUnlock()
		return c.conn, nil
	}

	// connection is nil, unlock read lock and lock write lock
	c.connLocker.RUnlock()
	c.connLocker.Lock()
	defer c.connLocker.Unlock()

	if c.conn != nil {
		// other goroutine got the lock first then set the connection, return it
		return c.conn, nil
	}
	c.conn, err = c.Connection.Connect()
	if err != nil {
		err = errs.New(err)
		return
	}
	ch := make(chan *amqp091.Error)
	go func() {
		e := <-ch
		fmt.Println("connection closed: ", e)

		// connection closed, set connection to nil
		c.connLocker.Lock()
		defer c.connLocker.Unlock()
		c.conn = nil
	}()
	c.conn.NotifyClose(ch)

	return c.conn, nil
}

type Consume struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp091.Table
}

func (c Consume) Consume(ch *amqp091.Channel, name string) (<-chan amqp091.Delivery, error) {
	return ch.Consume(name, c.Consumer, c.AutoAck, c.Exclusive, c.NoLocal, c.NoWait, c.Args)
}
