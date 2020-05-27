package main

import (
	log "github.com/jeanphorn/log4go"
	"github.com/programmer74/gotsdb/services"
	"github.com/programmer74/gotsdb/services/cluster"
	"github.com/programmer74/gotsdb/services/servers"
	"github.com/programmer74/gotsdb/services/storage/kvs"
	"github.com/programmer74/summer/summer"
)

const logSettingsFile = "./log4go.json"
const propertiesFile = "./app.properties"

const storageBeanName = "Storage"
const storageBeanType = "*kvs.Storage"
const applicationBeanName = "Application"
const grpcUserServerBeanName = "GrpcUserServer"
const grpcClusterServerBeanName = "GrpcClusterServer"
const clusterManagerBeanName = "ClusterManager"

const kvsEnginePropertyKey = "kvs.engine"
const kvsEnginePropertyFileValue = "file"
const kvsEnginePropertyInMemValue = "inmem"

func setupDI() {
	summer.ParseProperties(propertiesFile)

	summer.RegisterBean(applicationBeanName, services.Application{})
	summer.RegisterBean(grpcUserServerBeanName, servers.GrpcUserServer{})

	summer.RegisterBean(grpcClusterServerBeanName, cluster.GrpcClusterServer{})
	summer.RegisterBean(clusterManagerBeanName, cluster.Manager{})

	kvsEngine, _ := summer.GetPropertyValue(kvsEnginePropertyKey)
	switch kvsEngine {
	case kvsEnginePropertyFileValue:
		summer.RegisterBeanWithTypeAlias(storageBeanName, kvs.FileKVS{}, storageBeanType)
	case kvsEnginePropertyInMemValue:
		summer.RegisterBeanWithTypeAlias(storageBeanName, kvs.InMemKVS{}, storageBeanType)
	}

	summer.PerformDependencyInjection()
}

func main() {
	log.LoadConfiguration(logSettingsFile)
	defer log.Close()

	setupDI()

	app := summer.GetBean("Application").(*services.Application)
	app.Startup()
}