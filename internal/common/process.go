package common

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pastelnetwork/go-commons/log"
	"golang.org/x/sync/errgroup"
)

type Application struct {
	name string
}

func NewApplication(name string) *Application {
	return &Application{
		name: name,
	}
}

func (a *Application) Run(servers []func(ctx context.Context, a *Application) func() error) {

	log.Infof("")
	log.Infof("=======================================")
	log.Infof("====== %s starting ======", a.name)
	log.Infof("=======================================")

	if len(servers) == 0 {
		panic("runners array can't be empty")
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTERM)
		<-signals
		log.Infoln("program interrupted")
		cancel()
		log.Infoln("cancel context sent")
	}()

	eg, ctx := errgroup.WithContext(ctx)
	for _, f := range servers {
		eg.Go(f(ctx, a))
	}

	if err := eg.Wait(); err != nil {
		log.Errorf("error in the server goroutines: %s", err)
		os.Exit(1)
	}
	log.Infoln("everything closed successfully")
	log.Infoln("exiting")
}

func (a *Application) CreateServer(ctx context.Context, serverName string,
	startServer func(ctx context.Context) error,
	runServer func(ctx context.Context) error,
	stopServer func(ctx context.Context) error) func() error {

	return func() error {
		if err := startServer(ctx); err != nil {
			return fmt.Errorf("error starting the %s server: %w", serverName, err)
		}

		errChan := make(chan error, 1)

		go func() {
			<-ctx.Done()
			shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			log.Infof("stopping server %s", serverName)
			if err := stopServer(shutCtx); err != nil {
				errChan <- fmt.Errorf("error shutting down the %s server: %w", serverName, err)
			}
			log.Infof("the %s server is stopped", serverName)
			close(errChan)
		}()

		log.Infof("the %s server is starting", serverName)

		if err := runServer(ctx); err != nil {
			return fmt.Errorf("error running the %s server: %w", serverName, err)
		}
		log.Infof("the %s server is closing", serverName)

		err := <-errChan
		log.Infof("the %s server is closed", serverName)
		return err
	}
}
