package services

import (
	log "github.com/jeanphorn/log4go"
	"github.com/nikita-tomilov/gotsdb/services/cluster"
	"github.com/nikita-tomilov/gotsdb/services/servers"
	"github.com/nikita-tomilov/gotsdb/services/storage"
	"github.com/nikita-tomilov/summer/summer"
)

type Application struct {
	StorageManager  *interface{} `summer:"StorageManager"`
	GrpcUserServer *interface{} `summer:"GrpcUserServer"`
}

func (a *Application) getStorageManager() storage.Manager {
	s := *a.StorageManager
	s2 := (s).(storage.Manager)
	return s2
}

func (a *Application) getGrpcServer() *servers.GrpcUserServer {
	s := a.GrpcUserServer
	s2 := (*s).(*servers.GrpcUserServer)
	return s2
}

func (a *Application) getClusterManager() *cluster.Manager {
	s := summer.GetBean("ClusterManager")
	if s != nil {
		s2 := (*s).(*cluster.Manager)
		return s2
	}
	return nil
}

func (a *Application) Startup() {
	log.Warn("Autowiring OK")

	log.Warn("Initializing storage level...")
	s := a.getStorageManager()
	s.InitStorage()
	log.Warn("Storage level OK")

	c := a.getClusterManager()
	if c != nil {
		log.Warn("Launching Cluster Manager...")
		c.StartClustering()
		log.Warn("Cluster Manager startup OK")
	}

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
