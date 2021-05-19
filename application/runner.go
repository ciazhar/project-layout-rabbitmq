package application

import (
	"github.com/ciazhar/project-layout-rabbitmq/third_party"
	"github.com/ciazhar/project-layout-rabbitmq/third_party/logger"
	"os"
	"os/signal"
)

func AppRunner(daemon third_party.Daemon) error {
	err := daemon.Start()
	if err != nil {
		return err
	}

	osSignals := make(chan os.Signal)
	signal.Notify(osSignals, os.Interrupt)

	select {
	case <-osSignals:
		logger.Log.Infof("osSignal Interrupt trigerred")
		return daemon.Stop()
	}
}
