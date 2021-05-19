package amqp

import (
	"github.com/ciazhar/project-layout-rabbitmq/third_party/logger"
	"github.com/streadway/amqp"
)

type amqpClient struct {
	amqpPublisherManager
	amqpSubscriptionManager
	connection *amqp.Connection
}

func (cli *amqpClient) Init(connection *amqp.Connection) error {
	cli.connection = connection

	logger.Log.Infof("initiate publisher manager")
	if err := cli.amqpPublisherManager.Init(connection); err != nil {
		logger.Log.Warnf("Fail initiate publisher manager %v", err)
		return err
	}

	logger.Log.Infof("initiate subscription manager")
	if err := cli.amqpSubscriptionManager.Init(connection); err != nil {
		logger.Log.Warnf("Fail initiate subscription manager %v", err)
		return err
	}
	return nil
}

func (cli *amqpClient) Close() error {
	logger.Log.Infof("try close subscription manager")
	if err := cli.amqpSubscriptionManager.Close(); err != nil {
		logger.Log.Errorf("Failed to close subscription manager: %v\n", err)
	}

	logger.Log.Infof("try close publisher manager")
	if err := cli.amqpPublisherManager.Close(); err != nil {
		logger.Log.Errorf("Failed to close publisher manager: %v\n", err)
	}

	logger.Log.Infof("try close connection")
	return cli.connection.Close()
}
