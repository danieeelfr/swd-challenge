package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danieeelfr/swd-challenge/internal/config"
	httpServer "github.com/danieeelfr/swd-challenge/internal/httpserver"
	"github.com/danieeelfr/swd-challenge/pkg/wait"
	"github.com/sirupsen/logrus"
)

var (
	log            = logrus.WithField("package", "main.app")
	wg             = wait.New()
	waitToShutdown time.Duration
)

func init() {
	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = logrus.InfoLevel
	}

	logrus.SetLevel(level)
}

func main() {
	log.Info("starting api http server...")

	cfg := config.New(config.SwdApp)

	waitToShutdown = time.Duration(cfg.HTTPServerConfig.WaitToShutdown) * time.Second
	srv, err := httpServer.New(cfg, wg)
	if err != nil {
		log.WithError(err).Fatal(fmt.Sprintf("could not run %s", config.SwdApp))
	}

	wg.Add()
	shutdownSignal(srv)

	if err := srv.Start(); err != nil {
		log.WithError(err).Fatal(fmt.Sprintf("fail starting %s", config.SwdApp))
	}

	wg.Wait()

	log.Infof("finishing %s", config.SwdApp)

}

func shutdownSignal(ctrlHTTP httpServer.Interactor) {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-signalChannel
		switch sig {
		case syscall.SIGTERM, syscall.SIGINT:
			log.Infof("interruption request. signal: [%v].", sig)
			wg.Block()
			ctrlHTTP.Shutdown()

			log.Infof("waiting [%v] for open processes.", waitToShutdown)
			time.Sleep(waitToShutdown)

			log.Infof("finishing...")
			for wg.Done() {
				log.Infof("ignoring open process to kill...")
			}
		}
	}()
}
