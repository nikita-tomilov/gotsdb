package main

import (
	log "github.com/jeanphorn/log4go"
	"github.com/programmer74/gotsdb/services"
	"github.com/programmer74/gotsdb/services/cluster"
	"github.com/programmer74/gotsdb/services/servers"
	"github.com/programmer74/gotsdb/services/storage"
	"github.com/programmer74/gotsdb/services/storage/kvs"
	"github.com/programmer74/summer/summer"
	"os"
)

const logSettingsFile = "./log4go.json"
const defaultPropertiesFile = "./app.properties"

const propertiesOverrideEnvironmentVariable = "GOTSDB_PROPERTY_FILE"

const storageBeanName = "KeyValueStorage"
const storageBeanType = "*kvs.KeyValueStorage"
const applicationBeanName = "Application"
const grpcUserServerBeanName = "GrpcUserServer"
const grpcClusterServerBeanName = "GrpcClusterServer"
const clusterManagerBeanName = "ClusterManager"
const clusteredStorageManagerBeanName = "ClusteredStorageManager"

const kvsEnginePropertyKey = "kvs.engine"
const kvsEnginePropertyFileValue = "file"
const kvsEnginePropertyInMemValue = "inmem"

func setupDI() {
	propertiesOverride, propertiesOverridePresent:= os.LookupEnv(propertiesOverrideEnvironmentVariable)
	if propertiesOverridePresent {
		summer.ParseProperties(propertiesOverride)
	} else {
		summer.ParseProperties(defaultPropertiesFile)
	}

	summer.RegisterBean(applicationBeanName, services.Application{})
	summer.RegisterBean(grpcUserServerBeanName, servers.GrpcUserServer{})

	summer.RegisterBean(grpcClusterServerBeanName, cluster.GrpcClusterServer{})
	summer.RegisterBean(clusterManagerBeanName, cluster.Manager{})
	summer.RegisterBean(clusteredStorageManagerBeanName, storage.ClusteredStorageManager{})

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