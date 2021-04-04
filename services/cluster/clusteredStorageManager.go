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
	case "none":
		{
			c.readingConsistency = NONE
		}
	case "any":
		{
			c.readingConsistency = ANY
		}
	case "all":
		{
			c.readingConsistency = ALL
		}
	default:
		log.Critical("Unknown reading consistency %s", c.ReadingConsistencyStr)
	}
	switch c.WritingConsistencyStr {
	case "none":
		{
			c.writingConsistency = NONE
		}
	case "any":
		{
			c.writingConsistency = ANY
		}
	case "all":
		{
			c.writingConsistency = ALL
		}
	default:
		log.Critical("Unknown writing consistency %s", c.WritingConsistencyStr)
	}
	log.Warn("Reading consistency is set to '%s', writing consistency is set to '%s'", c.ReadingConsistencyStr, c.WritingConsistencyStr)
}

type MessageWithId interface {
	GetMsgId() uint32
}

func (c *ClusteredStorageManager) shouldSaveDataOnAllNodes() bool { return c.writingConsistency == ALL }

func (c *ClusteredStorageManager) proxyIfNeeded(id uint32, callback func(pb.ClusterClient) bool) {
	if !c.proxiedCommands.Contains(id) {
		c.proxiedCommands.Put(id)
		shouldContinue := true
		for _, o := range c.clusterManager.GetKnownOutboundConnections() {
			shouldContinue = callback(o.GetGrpcChannel())
			if !shouldContinue {
				break
			}
		}
	}
}

func (c *ClusteredStorageManager) KvsSave(ctx context.Context, req *pb.KvsStoreRequest) (*pb.KvsStoreResponse, error) {
	c.kvsStorage.Save(req.Key, req.Value)

	if c.shouldSaveDataOnAllNodes() {
		c.proxyIfNeeded(req.MsgId, func (otherNode pb.ClusterClient) bool {
			otherNode.KvsSave(ctx, req)
			return true
		})
	}

	return &pb.KvsStoreResponse{MsgId: req.MsgId, Ok: true}, nil
}

func (c *ClusteredStorageManager) shouldReadDataFromOtherNodes() bool { return c.readingConsistency != NONE }

func (c *ClusteredStorageManager) KvsKeyExists(ctx context.Context, req *pb.KvsKeyExistsRequest) (*pb.KvsKeyExistsResponse, error) {
	exists := c.kvsStorage.KeyExists(req.Key)

	if !exists && c.shouldReadDataFromOtherNodes() {
		c.proxyIfNeeded(req.MsgId, func (otherNode pb.ClusterClient) bool {
			existsOnAnotherNode, _ := otherNode.KvsKeyExists(ctx, req)
			if existsOnAnotherNode.Exists {
				exists = true
				return !exists
			}
			return true
		})
	}

	return &pb.KvsKeyExistsResponse{MsgId: req.MsgId, Exists: exists}, nil
}

func (c *ClusteredStorageManager) shouldSaveDataRetrievedFromOthers() bool { return c.writingConsistency == ANY }

func (c *ClusteredStorageManager) KvsRetrieve(ctx context.Context, req *pb.KvsRetrieveRequest) (*pb.KvsRetrieveResponse, error) {
	exists := c.kvsStorage.KeyExists(req.Key)
	value := make([]byte, 0)
	if exists {
		value = c.kvsStorage.Retrieve(req.Key)
		return &pb.KvsRetrieveResponse{MsgId: req.MsgId, Value: value}, nil
	}

	c.proxyIfNeeded(req.MsgId, func (otherNode pb.ClusterClient) bool {
		existsOnAnotherNode, _ := otherNode.KvsKeyExists(ctx, &pb.KvsKeyExistsRequest{MsgId: req.MsgId + 1, Key: req.Key})
		if existsOnAnotherNode.Exists {
			ans, err := otherNode.KvsRetrieve(ctx, &pb.KvsRetrieveRequest{MsgId: req.MsgId + 2, Key: req.Key})
			if err != nil {
				log.Error("Error in KvsRetrieve %s", err)
				return true
			}
			if c.shouldSaveDataRetrievedFromOthers() {
				c.kvsStorage.Save(req.Key, ans.Value)
			}
			return false
		}
		return true
	})

	return &pb.KvsRetrieveResponse{MsgId: req.MsgId, Value: value}, nil
}

func (c *ClusteredStorageManager) KvsDelete(ctx context.Context, req *pb.KvsDeleteRequest) (*pb.KvsDeleteResponse, error) {
	c.kvsStorage.Delete(req.Key)

	if c.shouldSaveDataOnAllNodes() {
		c.proxyIfNeeded(req.MsgId, func (otherNode pb.ClusterClient) bool {
			otherNode.KvsDelete(ctx, req)
			return true
		})
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


	if c.shouldReadDataFromOtherNodes() {
		c.proxyIfNeeded(req.MsgId, func (otherNode pb.ClusterClient) bool {
			knownKeysOnAnotherNode, err := otherNode.KvsGetKeys(ctx, &pb.KvsAllKeysRequest{MsgId: req.MsgId + 1})
			if err != nil {
				log.Error("Error in KvsGetKeys %s", err)
				return true
			}
			for _, knownKeyOnAnotherNode := range knownKeysOnAnotherNode.Keys {
				mapOfKnownKeys.Put(knownKeyOnAnotherNode)
			}
			return true
		})
	}

	var ans [][]byte
	for _, key := range mapOfKnownKeys.Values() {
		arr := key.([]byte)
		ans = append(ans, arr)
	}
	return &pb.KvsAllKeysResponse{MsgId: req.MsgId, Keys: ans}, nil
}

func (c *ClusteredStorageManager) TSSave(ctx context.Context, req *pb.TSStoreRequest) (*pb.TSStoreResponse, error) {
	c.tssStorage.Save(req.DataSource, req.Values, req.ExpirationMillis)

	if c.shouldSaveDataOnAllNodes() {
		c.proxyIfNeeded(req.MsgId,  func (otherNode pb.ClusterClient) bool {
			_, err := otherNode.TSSave(ctx, req)
			if err != nil {
				log.Error("Error in TSSave %s", err)
			}
			return true
		})
	}

	return &pb.TSStoreResponse{MsgId: req.MsgId, Ok: true}, nil
}

func (c *ClusteredStorageManager) TSSaveBatch(ctx context.Context, req *pb.TSStoreBatchRequest) (*pb.TSStoreResponse, error) {
	//TODO: IMPLEMENT
	panic("not implemented")
}

func mergeData(a map[string]*pb.TSPoints, b map[string]*pb.TSPoints) map[string]*pb.TSPoints {
	ans := make(map[string]*pb.TSPoints)
	for k, v := range a {
		existingData, exists := ans[k]
		if !exists {
			ans[k] = v
		} else {
			for ts, val := range v.Points {
				existingData.Points[ts] = val
			}
		}
	}
	for k, v := range b {
		existingData, exists := ans[k]
		if !exists {
			ans[k] = v
		} else {
			for ts, val := range v.Points {
				existingData.Points[ts] = val
			}
		}
	}
	return ans
}

func mergeAvail(a []*pb.TSAvailabilityChunk, b []*pb.TSAvailabilityChunk) []*pb.TSAvailabilityChunk {
	//TODO do something like Guava RangeSet
	min := (^uint64(0)) >> 1
	max := uint64(0)
	for _, c := range a {
		if min > c.FromTimestamp {
			min = c.FromTimestamp
		}
		if max < c.ToTimestamp {
			max = c.ToTimestamp
		}
	}
	for _, c := range b {
		if min > c.FromTimestamp {
			min = c.FromTimestamp
		}
		if max < c.ToTimestamp {
			max = c.ToTimestamp
		}
	}
	if min > max {
		return make([]*pb.TSAvailabilityChunk, 0)
	}
	return []*pb.TSAvailabilityChunk{{FromTimestamp: min, ToTimestamp: max}}
}

func (c *ClusteredStorageManager) shouldReadDataFromAllOtherNodes() bool { return c.readingConsistency == ALL }

func (c *ClusteredStorageManager) TSRetrieve(ctx context.Context, req *pb.TSRetrieveRequest) (*pb.TSRetrieveResponse, error) {
	ans := c.tssStorage.Retrieve(req.DataSource, req.Tags, req.FromTimestamp, req.ToTimestamp)

	if c.shouldReadDataFromOtherNodes() {
		c.proxyIfNeeded(req.MsgId, func (otherNode pb.ClusterClient) bool {
			ansFromOtherNode, err := otherNode.TSRetrieve(ctx, req)
			if err != nil {
				log.Error("Error in TSRetrieve %s", err)
				return true
			}
			ans = mergeData(ans, ansFromOtherNode.Values)
			if c.shouldSaveDataRetrievedFromOthers() {
				//TODO: FIX EXPIRATION, RETURN IT IN TSRetrieveResponse
				c.tssStorage.Save(ansFromOtherNode.DataSource, ansFromOtherNode.Values, 0)
			}
			return c.shouldReadDataFromAllOtherNodes()
		})
	}

	return &pb.TSRetrieveResponse{MsgId: req.MsgId, DataSource: req.DataSource, FromTimestamp: req.FromTimestamp, ToTimestamp: req.ToTimestamp, Values: ans}, nil
}

func (c *ClusteredStorageManager) TSAvailability(ctx context.Context, req *pb.TSAvailabilityRequest) (*pb.TSAvailabilityResponse, error) {
	ans := c.tssStorage.Availability(req.DataSource, req.FromTimestamp, req.ToTimestamp)

	if c.shouldReadDataFromOtherNodes() {
		c.proxyIfNeeded(req.MsgId, func (otherNode pb.ClusterClient) bool {
			ansFromOtherNode, err := otherNode.TSAvailability(ctx, req)
			if err != nil {
				log.Error("Error in TSAvailability %s", err)
				return true
			}
			ans = mergeAvail(ans, ansFromOtherNode.Availability)
			return c.shouldReadDataFromAllOtherNodes()
		})
	}

	return &pb.TSAvailabilityResponse{MsgId: req.MsgId, Availability: ans}, nil
}
