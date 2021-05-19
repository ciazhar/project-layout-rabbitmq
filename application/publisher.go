package application

import (
	"github.com/ciazhar/project-layout-rabbitmq/third_party"
	"github.com/ciazhar/project-layout-rabbitmq/third_party/amqp"
	logger2 "github.com/ciazhar/project-layout-rabbitmq/third_party/logger"
	messaging2 "github.com/ciazhar/project-layout-rabbitmq/third_party/messaging"
	"github.com/satori/go.uuid"
	"time"
)

type (
	publisher struct {
		queueName string
		broker    messaging2.Broker
		stopChan  chan bool
	}
)

func (app *Application) NewPublisherDaemon() third_party.Daemon {
	broker := amqp.NewAmqpBroker(app.rabbit)
	return &publisher{
		app.queueName,
		broker,
		make(chan bool),
	}
}

func (d *publisher) Start() error {
	err := d.broker.Start()
	if err != nil {
		return err
	}

	publisher, err := d.broker.CreatePublisher(d.queueName)
	if err != nil {
		return err
	}

	go d.runLoop(publisher)

	return nil
}

func (d *publisher) runLoop(publisher messaging2.Publisher) {
	logger := logger2.Log.WithField("contex", "publisher")
	for {
		select {
		default:
			logger.Debug("publishing started")
			publisher.Publish("hello "+uuid.NewV1().String(), uuid.NewV1().String())
			time.Sleep(time.Second * 1)
		case stop := <-d.stopChan:
			if stop {
				return
			}
		}
	}
}

func (d *publisher) Stop() error {
	d.stopChan <- true
	return d.broker.Stop()
}
