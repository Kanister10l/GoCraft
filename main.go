package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kanister10l/GoCraft/eventmanager"

	"github.com/kanister10l/GoCraft/logger"
)

func main() {
	fmt.Printf("%s", banner)
	logger.SetupLogger()
	logger.Logger.Infow("GoCraft is being started")

	eventmanager.InitMasterManager()
	eventmanager.Master.NewEvent("stop", func(param interface{}) interface{} {
		logger.Logger.Infof("Stop event executed")
		os.Exit(0)
		return nil
	})

	setupSignal()

	<-make(chan bool)
}

func setupSignal() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)

	go func() {
		signal := <-signalChannel
		logger.Logger.Info(fmt.Sprintf("Captured %s Signal", signal.String()))
		if signal.String() == "interrupt" {
			eventmanager.Master.ExecEvent("stop")
		}
	}()
}
