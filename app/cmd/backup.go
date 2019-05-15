// Copyright © 2019 Valentin Slyusarev <va.slyusarev@gmail.com>

package cmd

import (
	"compress/gzip"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/va-slyusarev/lgr"

	"github.com/va-slyusarev/pinky/app/config"
	"github.com/va-slyusarev/pinky/app/store"
)

type Backup struct {
	Cfg *config.BackupConfig
	DB  store.Storer
}

// Serv. Run backup service.
func (b *Backup) Serv(ctx context.Context) {
	if b.Cfg.Duration <= 0 {
		lgr.Warn("backup: service is not used: duration is %s", b.Cfg.Duration)
		return
	}
	if err := os.MkdirAll(b.Cfg.Path, os.ModePerm); err != nil {
		lgr.Error("backup: service is not used: path: %v", err)
		return
	}
	lgr.Info("backup: service is activated: path %q, duration %s", b.Cfg.Path, b.Cfg.Duration)
	tick := time.NewTicker(b.Cfg.Duration)
	lgr.Info("backup: first backup at %s", time.Now().Add(b.Cfg.Duration).Format("2006/01/02 15:04:05"))

	for {
		select {
		case <-tick.C:
			if err := b.create(); err != nil {
				lgr.Error("backup: %s", err)
				continue
			}

			b.remove()

			lgr.Info("backup: next backup at %s", time.Now().Add(b.Cfg.Duration).Format("2006/01/02 15:04:05"))
		case <-ctx.Done():
			lgr.Info("backup: service terminated")
			return
		}
	}
}

func (b *Backup) create() error {
	fileName := fmt.Sprintf("%s/backup-%s.gz", b.Cfg.Path, time.Now().Format("20060102T150405"))
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	gz := gzip.NewWriter(f)
	defer func() {
		_ = gz.Close()
	}()
	if err := b.DB.Backup(gz); err != nil {
		return err
	}
	lgr.Info("backup: create file %q", fileName)
	return nil
}

func (b *Backup) remove() {
	filesDir, err := ioutil.ReadDir(b.Cfg.Path)
	if err != nil {
		lgr.Error("remove: can't read files in path: %q: %v", b.Cfg.Path, err)
		return
	}
	var files []os.FileInfo
	for _, file := range filesDir {
		if strings.HasPrefix(file.Name(), "backup-") {
			files = append(files, file)
		}
	}
	// old -> new
	sort.Slice(files, func(i int, j int) bool { return files[i].Name() < files[j].Name() })

	if len(files) > b.Cfg.MaxCopies {
		for i := 0; i < len(files)-b.Cfg.MaxCopies; i++ {
			file := filepath.Join(b.Cfg.Path, files[i].Name())
			if err := os.Remove(file); err != nil {
				lgr.Error("backup: remove: can't delete file, skip: %q: %v", file, err)
				continue
			}
			lgr.Info("backup: remove: copy is outdated and has been deleted: %q", file)
		}
	}
}
