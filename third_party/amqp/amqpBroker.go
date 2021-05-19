package amqp

import (
	"github.com/ciazhar/project-layout-rabbitmq/third_party/logger"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

type (
	RabbitConfig struct {
		Host     string
		User     string
		Password string
	}

	AmqpBroker struct {
		sync.Mutex
		amqpClient
		url       string
		conn      *amqp.Connection
		watchStop chan bool
		errors    chan *amqp.Error
	}
)

func getRabbitConnectionString(rabbitConfig *RabbitConfig) string {
	credentials := ""

	if rabbitConfig.User != "" && rabbitConfig.Password != "" {
		credentials = rabbitConfig.User + ":" + rabbitConfig.Password + "@"
	}

	return "amqp://" + credentials + rabbitConfig.Host + "/"
}

func NewAmqpBroker(rabbitConfig *RabbitConfig) *AmqpBroker {
	url := getRabbitConnectionString(rabbitConfig)

	broker := AmqpBroker{
		sync.Mutex{},
		amqpClient{},
		url,
		nil,
		make(chan bool),
		nil,
	}

	return &broker
}

func (b *AmqpBroker) Start() error {
	err := b.setup()
	go b.watch()
	return err
}

func (b *AmqpBroker) Stop() error {
	b.stopWatch()
	defer close(b.watchStop)
	return b.Close()
}

func (b *AmqpBroker) setup() error {
	logger.Log.Info("Setup AMQP Connection")
	if b.conn != nil {
		return nil
	}

	conn, err := b.connect()
	if err != nil {
		logger.Log.Errorf("connection failed: %v", err)
		return err
	}

	err = b.Init(conn)
	if err != nil {
		logger.Log.Errorf("client init failed: %v", err)
		return err
	}

	b.Lock()
	b.conn = conn
	b.Unlock()

	errors := make(chan *amqp.Error)
	b.conn.NotifyClose(errors)
	b.errors = errors
	return nil
}

func (b *AmqpBroker) connect() (*amqp.Connection, error) {
	return amqp.Dial(b.url)
}

func (b *AmqpBroker) watch() {
	for {
		select {
		case err := <-b.errors:
			logger.Log.Warnf("Connection lost: %v\n", err)
			b.Lock()
			b.conn = nil
			b.Unlock()
			time.Sleep(10 * time.Second)
			b.setup()
		case stop := <-b.watchStop:
			if stop {
				return
			}
		}
	}
}

func (b *AmqpBroker) stopWatch() {
	b.watchStop <- true
}
