// Copyright © 2019 Valentin Slyusarev <va.slyusarev@gmail.com>

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/va-slyusarev/afs"
	"github.com/va-slyusarev/lgr"

	"github.com/va-slyusarev/pinky/app/config"
	"github.com/va-slyusarev/pinky/app/crypto"
	"github.com/va-slyusarev/pinky/app/store"
	"github.com/va-slyusarev/pinky/app/types"
)

type Server struct {
	Cfg  *config.ServerConfig
	FS   *afs.AFS
	DB   store.Storer
	H    crypto.Hasher
	http *http.Server
}

// Serv. Run rest service.
func (s *Server) Serv(ctx context.Context) {

	s.initServer(s.initHandler())

	errChan := make(chan error)
	go func() {
		errChan <- s.http.ListenAndServe()
	}()

	lgr.Warn(fmt.Sprintf("rest server: service started on port :%d", s.Cfg.RestPort))

	for {
		select {
		case <-ctx.Done():
			s.shutdown()
			return
		case err := <-errChan:
			lgr.Error("rest server: service terminated by error: %v", err)
			return
		}
	}
}

func (s *Server) initServer(handler http.Handler) {
	s.http = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.Cfg.RestPort),
		Handler:           handler,
		ReadHeaderTimeout: 15 * time.Second,
		WriteTimeout:      45 * time.Second,
		IdleTimeout:       45 * time.Second,
	}
}

func (s *Server) initHandler() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Throttle(1000))
	router.Use(middleware.Timeout(45 * time.Second))

	router.Get("/*", s.restoreHandler)
	router.Post("/", s.storeHandler)
	router.Get("/health/ping", s.pingHandler)

	return router
}

func (s *Server) restoreHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.RequestURI, "/")
	lgr.Debug("restore uri: request id: %q", id)

	if ok := s.FS.Belong(id); !ok {
		lgr.Debug("restore uri: id %q is not belong fs", id)
		item := &types.Item{ID: id}

		ok, err := s.DB.Get(item)
		if ok {
			lgr.Debug("restore uri: id %q found in store: %q", id, item.URI)
			http.Redirect(w, r, item.URI, http.StatusTemporaryRedirect)
			return
		}
		lgr.Error("restore uri: id: %q: store: %v", id, err)
		http.Redirect(w, r, "/404.html", http.StatusTemporaryRedirect)
		return
	}

	lgr.Debug("restore uri: id %q is belong fs", id)
	http.FileServer(s.FS).ServeHTTP(w, r)
}

func (s *Server) storeHandler(w http.ResponseWriter, r *http.Request) {
	buf, _ := ioutil.ReadAll(r.Body)
	_ = r.Body.Close()
	uri := string(buf)

	lgr.Debug("store uri: request uri: %q", uri)

	item := new(types.Item)
	item.SetURI(uri)
	item.SetID(s.H.Hash(item.URI))

	lgr.Debug("store uri: from item uri: %q", item.URI)

	if err := s.DB.Put(item); err != nil {
		lgr.Error("error store item: %v", err)
		return
	}
	lgr.Debug("store uri: success: id is %q", item.ID)

	_, _ = w.Write([]byte(item.ID))
}

func (s *Server) pingHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("pong"))
}

func (s *Server) shutdown() {
	lgr.Info("rest server: service shutdown...")
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	if err := s.http.Shutdown(ctx); err != nil {
		lgr.Error("rest server: shutdown service with error: %v", err)
	}
	lgr.Info("rest server: service terminated")
}
