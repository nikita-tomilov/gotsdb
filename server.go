package main

import (
	log "github.com/jeanphorn/log4go"
	"github.com/nikita-tomilov/gotsdb/services"
	"github.com/nikita-tomilov/gotsdb/services/cluster"
	"github.com/nikita-tomilov/gotsdb/services/servers"
	"github.com/nikita-tomilov/gotsdb/services/storage/kvs"
	"github.com/nikita-tomilov/summer/summer"
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
const clusteredStorageManagerBeanName = "ClusteredStorageManagerAutowired"

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
	summer.RegisterBean(clusteredStorageManagerBeanName, cluster.ClusteredStorageManager{})

	kvsEngine, _ := summer.GetPropertyValue(kvsEnginePropertyKey)
	switch kvsEngine {
	case kvsEnginePropertyFileValue:
		summer.RegisterBeanWithTypeAlias(storageBeanName, kvs.FileKVS{}, storageBeanType)
	case kvsEnginePropertyInMemValue:
		summer.RegisterBeanWithTypeAlias(storageBeanName, kvs.InMemKVS{}, storageBeanType)
	}

	summer.PerformDependencyInjection()

	summer.PrintDependencyGraphVertex()
}

func main() {
	log.LoadConfiguration(logSettingsFile)
	defer log.Close()

	setupDI()

	app := summer.GetBean("Application").(*services.Application)
	app.Startup()
}