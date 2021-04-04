package storage

import (
	"context"
	pb "github.com/nikita-tomilov/gotsdb/proto"
)

type Manager interface {
	InitStorage()
	KvsSave(ctx context.Context, req *pb.KvsStoreRequest) (*pb.KvsStoreResponse, error)
	KvsKeyExists(ctx context.Context, req *pb.KvsKeyExistsRequest) (*pb.KvsKeyExistsResponse, error)
	KvsRetrieve(ctx context.Context, req *pb.KvsRetrieveRequest) (*pb.KvsRetrieveResponse, error)
	KvsDelete(ctx context.Context, req *pb.KvsDeleteRequest) (*pb.KvsDeleteResponse, error)
	KvsGetKeys(ctx context.Context, req *pb.KvsAllKeysRequest) (*pb.KvsAllKeysResponse, error)
	TSSave(ctx context.Context, req *pb.TSStoreRequest) (*pb.TSStoreResponse, error)
	TSSaveBatch(ctx context.Context, req *pb.TSStoreBatchRequest) (*pb.TSStoreResponse, error)
	TSRetrieve(ctx context.Context, req *pb.TSRetrieveRequest) (*pb.TSRetrieveResponse, error)
	TSAvailability(ctx context.Context, req *pb.TSAvailabilityRequest) (*pb.TSAvailabilityResponse, error)
}
