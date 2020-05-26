package main

import (
	log "github.com/jeanphorn/log4go"
	"github.com/programmer74/gotsdb/storage/kvs"
)

type Application struct {
	KvsStorage *interface{} `summer:"*kvs.Storage"`
}

func (a *Application) getKvsStorage() *kvs.Storage {
	s := a.KvsStorage
	s2 := (*s).(*kvs.Storage)
	return s2
}

func (a *Application) startup() {
	log.Warn("Startup OK")
}
