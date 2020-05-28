package storage

import (
	"github.com/programmer74/gotsdb/services/cluster"
	"github.com/programmer74/gotsdb/services/storage/kvs"
)

type ClusteredStorageManager struct {
	KvsStorageAutowired *interface{} `summer:"*kvs.KeyValueStorage"`
	ClusterManagerAutowired *interface{} `summer:"*cluster.Manager"`
	kvsStorage kvs.KeyValueStorage
	clusterManager *cluster.Manager
}

func (c *ClusteredStorageManager) getKvsStorage() kvs.KeyValueStorage {
	ks := *c.KvsStorageAutowired
	ks2 := (ks).(kvs.KeyValueStorage)
	return ks2
}

func (c *ClusteredStorageManager) getClusterManager() *cluster.Manager {
	s := c.ClusterManagerAutowired
	s2 := (*s).(*cluster.Manager)
	return s2
}

func (c *ClusteredStorageManager) InitStorage() {
	c.kvsStorage = c.getKvsStorage()
	c.kvsStorage.InitStorage()
	c.clusterManager = c.getClusterManager()
}

func (c *ClusteredStorageManager) Save(key []byte, value []byte) {
	c.kvsStorage.Save(key, value)
}

func (c *ClusteredStorageManager) KeyExists(key []byte) bool {
	return c.kvsStorage.KeyExists(key)
}

func (c *ClusteredStorageManager) Retrieve(key []byte) []byte {
	return c.kvsStorage.Retrieve(key)
}

func (c *ClusteredStorageManager) Delete(key []byte) {
	c.kvsStorage.Delete(key)
}
