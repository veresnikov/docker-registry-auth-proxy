package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/veresnikov/docker-registry-auth-proxy/pkg/infrastructure/logger"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const (
	applicationID = "auth-proxy"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	ctx = listenOSKillSignalsContext(ctx)
	mainLogger := logger.NewJSONLogger(&logger.Config{
		AppName: applicationID,
	})
	cnf, err := parseEnv()
	if err != nil {
		mainLogger.FatalError(err, "failed to parse env")
	}
	app := &cli.App{
		Name: applicationID,
		Commands: cli.Commands{
			serve(cnf, mainLogger),
		},
	}
	err = app.RunContext(ctx, os.Args)
	switch errors.Cause(err) {
	case errServerIsStopped:
		mainLogger.Info(err)
	default:
		mainLogger.FatalError(err)
	}
}

func listenOSKillSignalsContext(ctx context.Context) context.Context {
	var cancelFunc context.CancelFunc
	ctx, cancelFunc = context.WithCancel(ctx)
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
		select {
		case <-ch:
			cancelFunc()
		case <-ctx.Done():
			return
		}
	}()
	return ctx
}
