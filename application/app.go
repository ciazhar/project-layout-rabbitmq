package application

import (
	"github.com/ciazhar/project-layout-rabbitmq/third_party/amqp"
	logger2 "github.com/ciazhar/project-layout-rabbitmq/third_party/logger"
)

type (
	Application struct {
		queueName string
		rabbit    *amqp.RabbitConfig
	}
)

func SetupApp() *Application {
	amqpConf := amqp.RabbitConfig{
		Host:     "localhost:5672",
		User:     "guest",
		Password: "guest",
	}

	logger := logger2.Util{Stdout: true, Level: "DEBUG"}
	logger2.Log = logger.NewLogger()

	return &Application{
		"testing",
		&amqpConf,
	}
}
