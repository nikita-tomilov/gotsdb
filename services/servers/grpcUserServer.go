package servers

import (
	log "github.com/jeanphorn/log4go"
	pb "github.com/programmer74/gotsdb/proto"
	"github.com/programmer74/gotsdb/services/storage/kvs"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
)

type GrpcUserServer struct{
	ListenAddress string `summer.property:"grpc.listenAddress|:5300"`
	KvsStorage *interface{} `summer:"*kvs.Storage"`
}

func (s *GrpcUserServer) getKvsStorage() kvs.Storage {
	ks := *s.KvsStorage
	ks2 := (ks).(kvs.Storage)
	return ks2
}

func (s *GrpcUserServer) BeginListening() {
	log.Warn("Starting to listen at '%s'", s.ListenAddress)
	listener, err := net.Listen("tcp", s.ListenAddress)

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterGoTSDBServer(grpcServer, &server{storage: s.getKvsStorage()})
	grpcServer.Serve(listener)
}

//TODO: pass errors from storage level to grpc level?

type server struct{
	storage kvs.Storage
}

func (s *server) KvsSave(c context.Context, req *pb.KvsStoreRequest) (*pb.KvsStoreResponse, error) {
	s.storage.Save(req.Key, req.Value)
	return &pb.KvsStoreResponse{Ok: true}, nil
}

func (s *server) KvsKeyExists(c context.Context, req *pb.KvsKeyExistsRequest) (*pb.KvsKeyExistsResponse, error) {
	exists := s.storage.KeyExists(req.Key)
	return &pb.KvsKeyExistsResponse{Exists: exists}, nil
}

func (s *server) KvsRetrieve(c context.Context, req *pb.KvsRetrieveRequest) (*pb.KvsRetrieveResponse, error) {
	value := s.storage.Retrieve(req.Key)
	return &pb.KvsRetrieveResponse{Value: value}, nil
}

func (s *server) KvsDelete(c context.Context, req *pb.KvsDeleteRequest) (*pb.KvsDeleteResponse, error) {
	s.storage.Delete(req.Key)
	return &pb.KvsDeleteResponse{Ok: true}, nil
}