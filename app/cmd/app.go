// Copyright © 2019 Valentin Slyusarev <va.slyusarev@gmail.com>

package cmd

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/va-slyusarev/afs"
	"github.com/va-slyusarev/lgr"

	"github.com/va-slyusarev/pinky/app/config"
	"github.com/va-slyusarev/pinky/app/crypto"
	"github.com/va-slyusarev/pinky/app/store"
	_ "github.com/va-slyusarev/pinky/web" // register assets
)

type application struct {
	cfg    *config.Config
	fs     *afs.AFS
	db     store.Storer
	h      crypto.Hasher
	ctx    context.Context
	cancel context.CancelFunc
	wait   chan struct{}
}

var app *application
var once sync.Once

// Run. Main entry point
func Run(ctx context.Context) error {

	newApp(ctx)

	// Init components
	app.initLogger()

	if err := app.initAFS(); err != nil {
		return err
	}

	if err := app.initStore(); err != nil {
		return err
	}

	if err := app.initHash(); err != nil {
		return err
	}

	// Start main loop
	go func() {
		select {
		case <-ctx.Done():
			app.cancel()
			time.Sleep(2 * time.Second) // shutdown services
			close(app.wait)
		}
	}()

	// Start services
	go app.srvBackup()
	go app.srvServer()

	app.Wait()
	lgr.Info("app: terminated")
	return nil
}

// newApp
func newApp(ctx context.Context) {
	once.Do(func() {
		app = new(application)
		app.cfg = &config.AppCfg
		app.ctx, app.cancel = context.WithCancel(ctx)
		app.wait = make(chan struct{})
	})
}

func (a *application) Wait() {
	<-a.wait
}

func (a *application) initLogger() {
	lgr.SetPrefix("pinky")
	lgr.SetLevel(a.cfg.Server.LogLevel)
	lgr.SetTpl(a.cfg.Server.LogTpl)
	lgr.Debug("logger initialized")
}

func (a *application) initAFS() error {
	fs, err := afs.GetAFS()
	if err != nil {
		return fmt.Errorf("initialized afs: %v", err)
	}
	// set dynamic params
	if err := fs.ExecTemplate([]string{"index.html", "404.html"},
		map[string]string{
			"domain": a.cfg.Server.Domain,
			"year":   time.Now().Format("2006"),
		}); err != nil {
		return fmt.Errorf("initialized afs: %v", err)
	}
	a.fs = fs
	lgr.Debug("afs initialized")
	return nil
}

func (a *application) initStore() error {
	db, err := store.Instance(a.cfg.Store)
	if err != nil {
		return fmt.Errorf("initialized store: %v", err)
	}
	a.db = db
	lgr.Debug("store initialized")
	return nil
}

func (a *application) initHash() error {
	h, err := crypto.Instance(a.cfg.Crypto)
	if err != nil {
		return fmt.Errorf("initialized crypto: %v", err)
	}
	a.h = h
	lgr.Debug("hash initialized")
	return nil
}

func (a *application) srvBackup() {
	b := &Backup{
		Cfg: a.cfg.Backup,
		DB:  a.db,
	}
	b.Serv(a.ctx)
}

func (a *application) srvServer() {
	s := &Server{
		Cfg: a.cfg.Server,
		FS:  a.fs,
		DB:  a.db,
		H:   a.h,
	}
	s.Serv(a.ctx)
}
