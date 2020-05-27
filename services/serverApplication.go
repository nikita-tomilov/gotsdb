package services

import (
	log "github.com/jeanphorn/log4go"
	"github.com/programmer74/gotsdb/services/cluster"
	"github.com/programmer74/gotsdb/services/servers"
	"github.com/programmer74/gotsdb/services/storage/kvs"
)

type Application struct {
	KvsStorage     *interface{} `summer:"*kvs.Storage"`
	GrpcUserServer *interface{} `summer:"*servers.GrpcUserServer"`
	ClusterManager *interface{} `summer:"*cluster.Manager"`
}

func (a *Application) getKvsStorage() kvs.Storage {
	s := *a.KvsStorage
	s2 := (s).(kvs.Storage)
	return s2
}

func (a *Application) getGrpcServer() *servers.GrpcUserServer {
	s := a.GrpcUserServer
	s2 := (*s).(*servers.GrpcUserServer)
	return s2
}

func (a *Application) getClusterManager() *cluster.Manager {
	s := a.ClusterManager
	s2 := (*s).(*cluster.Manager)
	return s2
}

func (a *Application) Startup() {
	log.Warn("Autowiring OK")

	log.Warn("Initializing storage level...")
	s := a.getKvsStorage()
	s.InitStorage()
	log.Warn("Storage level OK")

	log.Warn("Launching Cluster Manager...")
	c := a.getClusterManager()
	c.StartClustering()
	log.Warn("Cluster Manager startup OK")

	log.Warn("Launching GRPC User Server...")
	g := a.getGrpcServer()
	go g.BeginListening()
	log.Warn("GRPC User Server startup OK")

	a.blockMainThread()
}

func (a *Application) blockMainThread() {
	c := make(chan int)
	_ = <-c
}
