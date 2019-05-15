// Copyright © 2019 Valentin Slyusarev <va.slyusarev@gmail.com>

package config

import (
	"time"
)

var AppCfg Config

type Config struct {
	Server *ServerConfig `group:"server" namespace:"server" env-namespace:"SERVER"`
	Store  *StoreConfig  `group:"store" namespace:"store" env-namespace:"STORE"`
	Backup *BackupConfig `group:"backup" namespace:"backup" env-namespace:"BACKUP"`
	Crypto *CryptoConfig `group:"crypto" namespace:"crypto" env-namespace:"CRYPTO"`
}

type ServerConfig struct {
	LogLevel string `short:"l" long:"level" env:"LEVEL" default:"WARN" description:"logger level"`
	LogTpl   string `long:"tpl" env:"TPL" default:"md" description:"logger template"`
	Domain   string `short:"d" long:"domain" env:"DOMAIN" default:"http://localhost:8080" description:"domain name"`
	RestPort int    `short:"p" long:"port" env:"PORT" default:"8080" description:"rest api server port"`
}

type StoreConfig struct {
	DBPath  string        `long:"db" env:"DB" default:"pinky.db" description:"DB file"`
	Timeout time.Duration `long:"dbTimeout" env:"DB_TIMEOUT" default:"2s" description:"DB timeout"`
}

type BackupConfig struct {
	Path      string        `long:"path" env:"PATH" default:"./backup" description:"Backup path"`
	Duration  time.Duration `long:"duration" env:"DURATION" default:"0s" description:"Backup duration (0s - off)"`
	MaxCopies int           `long:"maxCopies" env:"MAX_COPIES" default:"7" description:"Backup max copies"`
}

type CryptoConfig struct {
	Salt uint64 `long:"salt" env:"SALT" default:"0x9876543210" description:"Hash salt"`
}
