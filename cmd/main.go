package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-template/cmd/option"
	"go-template/config"
	"go-template/internal/rest"
	"go-template/pkg/env"
	"go-template/pkg/json"
	"go-template/pkg/zlog"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = env.AppName()
	app.Usage = env.AppName() + " running"
	app.Commands = []cli.Command{
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "check version",
			Action: func(c *cli.Context) error {
				envInfo := env.CompileInfo()
				s := json.ToJSONs(envInfo)
				fmt.Println(s)
				return nil
			},
		},
	}
	app.Action = func(c *cli.Context) {
		appStart(c)
	}
	err := app.Run(os.Args)
	if err != nil {
		zlog.Fatalf(err.Error())
	}
}

func appStart(c *cli.Context) {
	opts := option.New()
	conf, err := opts.Config(c.String("config"))
	if err != nil {
		zlog.Fatalf(err.Error())
	}
	if err := config.DumpAndSetGlobal(conf); err != nil {
		zlog.Error("config_dump", zlog.Any("err", err.Error()))
	}
	zlog.Info("app_running", zlog.Any("env", env.CompileInfo()))
	startModules()
	gracefulShutdown()
	_ = zlog.Sync()
}

func startModules() {
	rest.Start()
}

func gracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT)
	s := <-c
	zlog.Info("signal", zlog.String("signal", s.String()))
}
