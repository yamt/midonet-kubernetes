package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/yamt/midonet-kubernetes/pkg/config"
	"github.com/yamt/midonet-kubernetes/pkg/factory"
	"github.com/yamt/midonet-kubernetes/pkg/midonet"
	"github.com/yamt/midonet-kubernetes/pkg/util"
)

func main() {
	app := cli.NewApp()
	app.Name = "midonetkube"
	app.Usage = "A Kubernetes controller for MidoNet integration"
	app.Version = "1.0.0"
	app.Action = RunController

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func RunController(ctx *cli.Context) error {
	if err := config.InitConfig(ctx, nil); err != nil {
		return err
	}

	clientset, err := util.NewClientset(&config.Kubernetes)
	if err != nil {
		return err
	}

	stopChan := make(chan struct{})
	factory, err := factory.NewWatchFactory(clientset, stopChan)
	if err != nil {
		return err
	}

	controller := midonet.NewMidoNetController(clientset, factory)
	controller.Run()

	select {}

	return nil
}
