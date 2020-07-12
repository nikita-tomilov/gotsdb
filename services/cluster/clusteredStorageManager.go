package cluster

import (
	"context"
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/services/storage"
	"github.com/nikita-tomilov/gotsdb/services/storage/kvs"
	"github.com/nikita-tomilov/gotsdb/services/storage/tss"
)

type ClusteredStorageManager struct {
	KvsStorageAutowired *interface{} `summer:"*kvs.KeyValueStorage"`
	kvsStorage          kvs.KeyValueStorage
	TssStorageAutowired *interface{} `summer:"*tss.TimeSeriesStorage"`
	tssStorage          tss.TimeSeriesStorage
	clusterManager      *Manager
	proxiedCommands     *storage.TTLSet
}

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
	c.kvsStorage = c.getKvsStorage()
	c.kvsStorage.InitStorage()
	c.tssStorage = c.getTsStorage()
	c.tssStorage.InitStorage()
	c.proxiedCommands = storage.NewTTLSet(10000, 10)
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

	if !exists && !c.proxiedCommands.Contains(req.MsgId) {
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

	if !exists && !c.proxiedCommands.Contains(req.MsgId) {
		c.proxiedCommands.Put(req.MsgId)
		for _, o := range c.clusterManager.GetKnownOutboundConnections() {
			existsOnAnotherNode, _ := o.GetGrpcChannel().KvsKeyExists(ctx, &pb.KvsKeyExistsRequest{MsgId: req.MsgId + 1, Key: req.Key})
			if existsOnAnotherNode.Exists {
				ans, err := o.GetGrpcChannel().KvsRetrieve(ctx, &pb.KvsRetrieveRequest{MsgId: req.MsgId + 1, Key: req.Key})
				if err == nil {
					c.kvsStorage.Save(req.Key, ans.Value)
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

	if !c.proxiedCommands.Contains(req.MsgId) {
		c.proxiedCommands.Put(req.MsgId)
		for _, o := range c.clusterManager.GetKnownOutboundConnections() {
			o.GetGrpcChannel().KvsDelete(ctx, req)
		}
	}

	return &pb.KvsDeleteResponse{MsgId: req.MsgId, Ok: true}, nil
}

func (c *ClusteredStorageManager) TSSave(ctx context.Context, req *pb.TSStoreRequest) (*pb.TSStoreResponse, error) {
	c.tssStorage.Save(req.DataSource, req.Values, req.ExpirationMillis)
	return &pb.TSStoreResponse{MsgId: req.MsgId, Ok: true}, nil
}

func (c *ClusteredStorageManager) TSRetrieve(ctx context.Context, req *pb.TSRetrieveRequest) (*pb.TSRetrieveResponse, error) {
	ans := c.tssStorage.Retrieve(req.DataSource, req.Tags, req.FromTimestamp, req.ToTimestamp)
	return &pb.TSRetrieveResponse{MsgId: req.MsgId, DataSource: req.DataSource, FromTimestamp: req.FromTimestamp, ToTimestamp: req.ToTimestamp, Values: ans}, nil
}
