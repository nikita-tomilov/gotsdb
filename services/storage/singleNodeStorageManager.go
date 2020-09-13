package storage

import (
	"bytes"
	"context"
	log "github.com/jeanphorn/log4go"
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/services/storage/kvs"
	"github.com/nikita-tomilov/gotsdb/services/storage/tss"
	"github.com/nikita-tomilov/gotsdb/utils"
)

type SingleNodeStorageManager struct {
	KvsStorageAutowired *interface{} `summer:"*kvs.KeyValueStorage"`
	kvsStorage          kvs.KeyValueStorage
	TssStorageAutowired *interface{} `summer:"*tss.TimeSeriesStorage"`
	tssStorage          tss.TimeSeriesStorage
}

func (c *SingleNodeStorageManager) getKvsStorage() kvs.KeyValueStorage {
	ks := *c.KvsStorageAutowired
	ks2 := (ks).(kvs.KeyValueStorage)
	return ks2
}

func (c *SingleNodeStorageManager) getTsStorage() tss.TimeSeriesStorage {
	ts := *c.TssStorageAutowired
	ts2 := (ts).(tss.TimeSeriesStorage)
	return ts2
}

func (c *SingleNodeStorageManager) InitStorage() {
	log.Warn("Using SingleNode StorageManager")
	c.kvsStorage = c.getKvsStorage()
	c.kvsStorage.InitStorage()
	log.Warn("Using '%s' as KeyValue storage backend", c.kvsStorage.String())
	c.tssStorage = c.getTsStorage()
	c.tssStorage.InitStorage()
	log.Warn("Using '%s' as TimeSeries storage backend", c.tssStorage.String())
}

func (c *SingleNodeStorageManager) KvsSave(ctx context.Context, req *pb.KvsStoreRequest) (*pb.KvsStoreResponse, error) {
	c.kvsStorage.Save(req.Key, req.Value)
	return &pb.KvsStoreResponse{MsgId: req.MsgId, Ok: true}, nil
}

func (c *SingleNodeStorageManager) KvsKeyExists(ctx context.Context, req *pb.KvsKeyExistsRequest) (*pb.KvsKeyExistsResponse, error) {
	exists := c.kvsStorage.KeyExists(req.Key)
	return &pb.KvsKeyExistsResponse{MsgId: req.MsgId, Exists: exists}, nil
}

func (c *SingleNodeStorageManager) KvsRetrieve(ctx context.Context, req *pb.KvsRetrieveRequest) (*pb.KvsRetrieveResponse, error) {
	value := c.kvsStorage.Retrieve(req.Key)
	return &pb.KvsRetrieveResponse{MsgId: req.MsgId, Value: value}, nil
}

func (c *SingleNodeStorageManager) KvsDelete(ctx context.Context, req *pb.KvsDeleteRequest) (*pb.KvsDeleteResponse, error) {
	c.kvsStorage.Delete(req.Key)
	return &pb.KvsDeleteResponse{MsgId: req.MsgId, Ok: true}, nil
}

func (c *SingleNodeStorageManager) KvsGetKeys(ctx context.Context, req *pb.KvsAllKeysRequest) (*pb.KvsAllKeysResponse, error) {

	mapOfKnownKeys := utils.NewHashSet(func(i interface{}) uint32 {
		b := i.([]byte)
		return utils.ComputeHashCode(b)
	}, func(a, b interface{}) bool {
		a2 := a.([]byte)
		b2 := b.([]byte)
		return bytes.Compare(a2, b2) == 0
	})

	knownKeysOnThisNode := c.kvsStorage.GetAllKeys()
	for _, knownKeyOnThisNode := range knownKeysOnThisNode {
		mapOfKnownKeys.Put(knownKeyOnThisNode)
	}

	var ans [][]byte
	for _, key := range mapOfKnownKeys.Values() {
		arr := key.([]byte)
		ans = append(ans, arr)
	}
	return &pb.KvsAllKeysResponse{MsgId: req.MsgId, Keys: ans}, nil
}

func (c *SingleNodeStorageManager) TSSave(ctx context.Context, req *pb.TSStoreRequest) (*pb.TSStoreResponse, error) {
	c.tssStorage.Save(req.DataSource, req.Values, req.ExpirationMillis)
	return &pb.TSStoreResponse{MsgId: req.MsgId, Ok: true}, nil
}

func (c *SingleNodeStorageManager) TSRetrieve(ctx context.Context, req *pb.TSRetrieveRequest) (*pb.TSRetrieveResponse, error) {
	ans := c.tssStorage.Retrieve(req.DataSource, req.Tags, req.FromTimestamp, req.ToTimestamp)
	return &pb.TSRetrieveResponse{MsgId: req.MsgId, DataSource: req.DataSource, FromTimestamp: req.FromTimestamp, ToTimestamp: req.ToTimestamp, Values: ans}, nil
}

func (c *SingleNodeStorageManager) TSAvailability(ctx context.Context, req *pb.TSAvailabilityRequest) (*pb.TSAvailabilityResponse, error) {
	ans := c.tssStorage.Availability(req.DataSource, req.FromTimestamp, req.ToTimestamp)
	return &pb.TSAvailabilityResponse{MsgId: req.MsgId, Availability: ans}, nil
}
