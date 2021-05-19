package main

import (
	"github.com/ciazhar/project-layout-rabbitmq/application"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := application.SetupApp()
	daemon := app.NewSubscriberDaemon()

	clientApp := cli.NewApp()
	clientApp.Name = "publisher-app"
	clientApp.Version = "0.0.1"
	clientApp.HideVersion = true
	clientApp.HideHelp = true
	clientApp.Action = application.AppRunner(daemon)
	clientApp.Run(os.Args)
}
