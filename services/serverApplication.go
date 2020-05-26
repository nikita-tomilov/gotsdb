package services

import (
	log "github.com/jeanphorn/log4go"
	"github.com/programmer74/gotsdb/services/servers"
	"github.com/programmer74/gotsdb/services/storage/kvs"
)

type Application struct {
	KvsStorage *interface{} `summer:"*kvs.Storage"`
	GrpcServer *interface{} `summer:"*servers.GrpcServer"`
}

func (a *Application) getKvsStorage() kvs.Storage {
	s := *a.KvsStorage
	s2 := (s).(kvs.Storage)
	return s2
}

func (a *Application) getGrpcServer() *servers.GrpcServer {
	s := a.GrpcServer
	s2 := (*s).(*servers.GrpcServer)
	return s2
}

func (a *Application) Startup() {
	log.Warn("Autowiring OK")

	log.Warn("Initializing storage level...")
	s := a.getKvsStorage()
	s.InitStorage()
	log.Warn("Storage level OK")


	log.Warn("Launching GRPC Server...")
	g := a.getGrpcServer()
	go g.Start()
	log.Warn("GRPC Server startup OK")

	a.blockMainThread()
}

func (a *Application) blockMainThread() {
	c := make(chan int)
	_ = <-c
}