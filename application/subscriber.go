package application

import (
	"context"
	"github.com/ciazhar/project-layout-rabbitmq/third_party"
	"github.com/ciazhar/project-layout-rabbitmq/third_party/amqp"
	messaging2 "github.com/ciazhar/project-layout-rabbitmq/third_party/messaging"
	"time"
)

type (
	subscriber struct {
		queueName string
		broker    messaging2.Broker
	}
)

func (app *Application) NewSubscriberDaemon() third_party.Daemon {
	broker := amqp.NewAmqpBroker(app.rabbit)
	return &subscriber{
		app.queueName,
		broker,
	}
}

func (sub *subscriber) Start() error {
	err := sub.broker.Start()
	if err != nil {
		return err
	}

	_, err = sub.broker.CreateSubscription(sub.queueName, sub.queueName, "", true, 5, sub.handle)
	if err != nil {
		return err
	}

	return nil
}

func (sub *subscriber) Stop() error {
	return sub.broker.Stop()
}

func (sub *subscriber) handle(ctx context.Context, event messaging2.Event) error {
	third_party.SessionLogger(ctx).Debugf("received : %s", string(event.GetBody()))
	time.Sleep(time.Second * 5)
	return nil
}
