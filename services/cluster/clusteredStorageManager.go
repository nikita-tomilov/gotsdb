package cluster

import (
	"bytes"
	"context"
	log "github.com/jeanphorn/log4go"
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/services/storage"
	"github.com/nikita-tomilov/gotsdb/services/storage/kvs"
	"github.com/nikita-tomilov/gotsdb/services/storage/tss"
	"github.com/nikita-tomilov/gotsdb/utils"
)

type ClusteredStorageManager struct {
	KvsStorageAutowired   *interface{} `summer:"*kvs.KeyValueStorage"`
	kvsStorage            kvs.KeyValueStorage
	TssStorageAutowired   *interface{} `summer:"*tss.TimeSeriesStorage"`
	tssStorage            tss.TimeSeriesStorage
	clusterManager        *Manager
	proxiedCommands       *storage.TTLSet
	ReadingConsistencyStr string `summer.property:"*cluster.readingConsistency|none"`
	WritingConsistencyStr string `summer.property:"*cluster.writingConsistency|none"`
	readingConsistency    int
	writingConsistency    int
}

const NONE = 0
const ANY = 1
const ALL = 2

func (c *ClusteredStorageManager) getKvsStorage() kvs.KeyValueStorage {
	ks := *c.KvsStorageAutowired
	ks2 := (ks).(kvs.KeyValueStorage)
	return ks2
}

func (c *ClusteredStorageManager) getTsStorage() tss.TimeSeriesStorage {
	ts := *c.TssStorageAutowired
	ts2 := (ts).(tss.TimeSeriesStorage)
	return ts2
}

func (c *ClusteredStorageManager) InitStorage() {
	log.Warn("Using ClusteredStorageManager")
	c.kvsStorage = c.getKvsStorage()
	c.kvsStorage.InitStorage()
	log.Warn("Using '%s' as KeyValue storage backend", c.kvsStorage.String())
	c.tssStorage = c.getTsStorage()
	c.tssStorage.InitStorage()
	log.Warn("Using '%s' as TimeSeries storage backend", c.tssStorage.String())
	c.parseConsistencies()
	c.proxiedCommands = storage.NewTTLSet(10000, 10)
}

func (c *ClusteredStorageManager) parseConsistencies() {
	switch c.ReadingConsistencyStr {
	case "none": {
		c.readingConsistency = NONE
	}
	case "any": {
		c.readingConsistency = ANY
	}
	case "all": {
		c.readingConsistency = ALL
	}
	default:
		log.Critical("Unknown reading consistency %s", c.ReadingConsistencyStr)
	}
	switch c.WritingConsistencyStr {
	case "none": {
		c.writingConsistency = NONE
	}
	case "any": {
		c.writingConsistency = ANY
	}
	case "all": {
		c.writingConsistency = ALL
	}
	default:
		log.Critical("Unknown writing consistency %s", c.WritingConsistencyStr)
	}
	log.Warn("Reading consistency is set to '%s', writing consistency is set to '%s'", c.ReadingConsistencyStr, c.WritingConsistencyStr)
}

func (c *ClusteredStorageManager) KvsSave(ctx context.Context, req *pb.KvsStoreRequest) (*pb.KvsStoreResponse, error) {
	c.kvsStorage.Save(req.Key, req.Value)

	if !c.proxiedCommands.Contains(req.MsgId) {
		c.proxiedCommands.Put(req.MsgId)
		for _, o := range c.clusterManager.GetKnownOutboundConnections() {
			o.GetGrpcChannel().KvsSave(ctx, req)
		}
	}

	return &pb.KvsStoreResponse{MsgId: req.MsgId, Ok: true}, nil
}

func (c *ClusteredStorageManager) KvsKeyExists(ctx context.Context, req *pb.KvsKeyExistsRequest) (*pb.KvsKeyExistsResponse, error) {
	exists := c.kvsStorage.KeyExists(req.Key)

	if !exists && !c.proxiedCommands.Contains(req.MsgId) && (c.readingConsistency != NONE) {
		c.proxiedCommands.Put(req.MsgId)
		for _, o := range c.clusterManager.GetKnownOutboundConnections() {
			existsOnAnotherNode, _ := o.GetGrpcChannel().KvsKeyExists(ctx, req)
			if existsOnAnotherNode.Exists {
				exists = true
				break
			}
		}
	}

	return &pb.KvsKeyExistsResponse{MsgId: req.MsgId, Exists: exists}, nil
}

func (c *ClusteredStorageManager) KvsRetrieve(ctx context.Context, req *pb.KvsRetrieveRequest) (*pb.KvsRetrieveResponse, error) {
	exists := c.kvsStorage.KeyExists(req.Key)
	value := c.kvsStorage.Retrieve(req.Key)
	//TODO: refactor this
	if !exists && !c.proxiedCommands.Contains(req.MsgId) && (c.readingConsistency != NONE) {
		c.proxiedCommands.Put(req.MsgId)
		for _, o := range c.clusterManager.GetKnownOutboundConnections() {
			existsOnAnotherNode, _ := o.GetGrpcChannel().KvsKeyExists(ctx, &pb.KvsKeyExistsRequest{MsgId: req.MsgId + 1, Key: req.Key})
			if existsOnAnotherNode.Exists {
				ans, err := o.GetGrpcChannel().KvsRetrieve(ctx, &pb.KvsRetrieveRequest{MsgId: req.MsgId + 1, Key: req.Key})
				if err == nil {
					if c.readingConsistency == ALL {
						c.kvsStorage.Save(req.Key, ans.Value)
					}
					return &pb.KvsRetrieveResponse{MsgId: req.MsgId, Value: ans.Value}, nil
				}
				return &pb.KvsRetrieveResponse{MsgId: req.MsgId, Value: value}, err
			}
		}
	}

	return &pb.KvsRetrieveResponse{MsgId: req.MsgId, Value: value}, nil
}

func (c *ClusteredStorageManager) KvsDelete(ctx context.Context, req *pb.KvsDeleteRequest) (*pb.KvsDeleteResponse, error) {
	c.kvsStorage.Delete(req.Key)

	if !c.proxiedCommands.Contains(req.MsgId) && (c.readingConsistency != NONE) {
		c.proxiedCommands.Put(req.MsgId)
		for _, o := range c.clusterManager.GetKnownOutboundConnections() {
			o.GetGrpcChannel().KvsDelete(ctx, req)
		}
	}

	return &pb.KvsDeleteResponse{MsgId: req.MsgId, Ok: true}, nil
}

func (c *ClusteredStorageManager) KvsGetKeys(ctx context.Context, req *pb.KvsAllKeysRequest) (*pb.KvsAllKeysResponse, error) {

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

	if !c.proxiedCommands.Contains(req.MsgId) && (c.readingConsistency != NONE) {
		c.proxiedCommands.Put(req.MsgId)
		for _, o := range c.clusterManager.GetKnownOutboundConnections() {
			knownKeysOnAnotherNode, err := o.GetGrpcChannel().KvsGetKeys(ctx, req)
			if err == nil {
				for _, knownKeyOnAnotherNode := range knownKeysOnAnotherNode.Keys {
					mapOfKnownKeys.Put(knownKeyOnAnotherNode)
				}
			}
		}
	}

	var ans [][]byte
	for _, key := range mapOfKnownKeys.Values() {
		arr := key.([]byte)
		ans = append(ans, arr)
	}
	return &pb.KvsAllKeysResponse{MsgId: req.MsgId, Keys: ans}, nil
}

func (c *ClusteredStorageManager) TSSave(ctx context.Context, req *pb.TSStoreRequest) (*pb.TSStoreResponse, error) {
	//TODO: cluster mode
	c.tssStorage.Save(req.DataSource, req.Values, req.ExpirationMillis)
	return &pb.TSStoreResponse{MsgId: req.MsgId, Ok: true}, nil
}

func (c *ClusteredStorageManager) TSRetrieve(ctx context.Context, req *pb.TSRetrieveRequest) (*pb.TSRetrieveResponse, error) {
	//TODO: cluster mode
	ans := c.tssStorage.Retrieve(req.DataSource, req.Tags, req.FromTimestamp, req.ToTimestamp)
	return &pb.TSRetrieveResponse{MsgId: req.MsgId, DataSource: req.DataSource, FromTimestamp: req.FromTimestamp, ToTimestamp: req.ToTimestamp, Values: ans}, nil
}

func (c *ClusteredStorageManager) TSAvailability(ctx context.Context, req *pb.TSAvailabilityRequest) (*pb.TSAvailabilityResponse, error) {
	//TODO: cluster mode
	ans := c.tssStorage.Availability(req.DataSource, req.FromTimestamp, req.ToTimestamp)
	return &pb.TSAvailabilityResponse{MsgId: req.MsgId, Availability: ans}, nil
}
