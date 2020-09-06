package main

import (
	log "github.com/jeanphorn/log4go"
	"github.com/nikita-tomilov/gotsdb/services"
	"github.com/nikita-tomilov/gotsdb/services/cluster"
	"github.com/nikita-tomilov/gotsdb/services/servers"
	"github.com/nikita-tomilov/gotsdb/services/storage/kvs"
	"github.com/nikita-tomilov/gotsdb/services/storage/tss"
	"github.com/nikita-tomilov/summer/summer"
	"os"
)

const logSettingsFile = "./log4go.json"
const defaultPropertiesFile = "./app.properties"

const propertiesOverrideEnvironmentVariable = "GOTSDB_PROPERTY_FILE"

const kvsStorageBeanName = "KeyValueStorage"
const kvsStorageBeanType = "*kvs.KeyValueStorage"
const tssStorageBeanName = "TimeSeriesStorage"
const tssStorageBeanType = "*tss.TimeSeriesStorage"
const applicationBeanName = "Application"
const grpcUserServerBeanName = "GrpcUserServer"
const grpcClusterServerBeanName = "GrpcClusterServer"
const clusterManagerBeanName = "ClusterManager"
const clusteredStorageManagerBeanName = "ClusteredStorageManagerAutowired"

const kvsEnginePropertyKey = "kvs.engine"
const kvsEnginePropertyFileValue = "file"
const kvsEnginePropertyInMemValue = "inmem"

const tssEnginePropertyKey = "tss.engine"
const tssEnginePropertyFileValue = "file"
const tssEnginePropertyInMemValue = "inmem"
const tssEnginePropertyLSMValue = "lsm"

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
		summer.RegisterBeanWithTypeAlias(kvsStorageBeanName, kvs.FileKVS{}, kvsStorageBeanType)
	case kvsEnginePropertyInMemValue:
		summer.RegisterBeanWithTypeAlias(kvsStorageBeanName, kvs.InMemKVS{}, kvsStorageBeanType)
	}

	tssEngine, _ := summer.GetPropertyValue(tssEnginePropertyKey)
	switch tssEngine {
	case tssEnginePropertyInMemValue:
		summer.RegisterBeanWithTypeAlias(tssStorageBeanName, tss.InMemTSS{}, tssStorageBeanType)
	case tssEnginePropertyFileValue:
		summer.RegisterBeanWithTypeAlias(tssStorageBeanName, tss.QlBasedPersistentTSS{}, tssStorageBeanType)
	case tssEnginePropertyLSMValue:
		summer.RegisterBeanWithTypeAlias(tssStorageBeanName, tss.LSMTSS{}, tssStorageBeanType)
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