// Copyright © 2019 Valentin Slyusarev <va.slyusarev@gmail.com>

//go:generate afs -src=web/asset -rw
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/va-slyusarev/lgr"

	"github.com/va-slyusarev/pinky/app/cmd"
	"github.com/va-slyusarev/pinky/app/config"
)

var revision = "develop"

func main() {
	fmt.Printf("pinky: app: revision: %s\n", revision)

	var p = flags.NewParser(&config.AppCfg, flags.Default)
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-stop
		lgr.Warn("app: detect interrupt, terminate after few seconds...")
		cancel()
	}()

	if err := cmd.Run(ctx); err != nil {
		lgr.Error("app: terminated: %v", err)
	}
}
