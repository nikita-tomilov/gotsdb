package cluster

import (
	log "github.com/jeanphorn/log4go"
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
)

type GrpcClusterServer struct{
	ListenAddress                    string       `summer.property:"cluster.listenAddress|:5300"`
	ClusteredStorageManagerAutowired *interface{} `summer:"*cluster.ClusteredStorageManager"`
	manager                          *Manager
	storageManager                   *ClusteredStorageManager
}

func (s *GrpcClusterServer) getStorageManager() *ClusteredStorageManager {
	sm := *s.ClusteredStorageManagerAutowired
	sm2 := (sm).(*ClusteredStorageManager)
	return sm2
}

func (s *GrpcClusterServer) Start() {
	log.Warn("Starting to listen for other nodes at '%s'", s.ListenAddress)
	listener, err := net.Listen("tcp", s.ListenAddress)

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	s.storageManager = s.getStorageManager()
	s.storageManager.clusterManager = s.manager
	pb.RegisterClusterServer(grpcServer, &clusterServer{parent:s})
	go grpcServer.Serve(listener)
}

//TODO: pass errors from storage level to grpc level?

type clusterServer struct{
	parent *GrpcClusterServer
}

func (s *clusterServer) Hello(c context.Context, rq *pb.HelloRequest) (*pb.AliveNodesResponse, error) {
	s.parent.manager.AddKnownNode(rq.Iam)
	return &pb.AliveNodesResponse{AliveNodes:s.parent.manager.GetKnownNodes()}, nil
}

func (s *clusterServer) Ping(c context.Context, rq *pb.PingRequest) (*pb.PingResponse, error) {
	log.Debug("Got ping rq")
	return &pb.PingResponse{Payload:rq.Payload}, nil
}

func (s *clusterServer) GetAliveNodes(c context.Context, v *pb.Void) (*pb.AliveNodesResponse, error) {
	return &pb.AliveNodesResponse{AliveNodes:s.parent.manager.GetKnownNodes()}, nil
}

func (s *clusterServer) KvsSave(c context.Context, req *pb.KvsStoreRequest) (*pb.KvsStoreResponse, error) {
	return s.parent.storageManager.KvsSave(c, req)
}

func (s *clusterServer) KvsKeyExists(c context.Context, req *pb.KvsKeyExistsRequest) (*pb.KvsKeyExistsResponse, error) {
	return s.parent.storageManager.KvsKeyExists(c, req)
}

func (s *clusterServer) KvsRetrieve(c context.Context, req *pb.KvsRetrieveRequest) (*pb.KvsRetrieveResponse, error) {
	return s.parent.storageManager.KvsRetrieve(c, req)
}

func (s *clusterServer) KvsDelete(c context.Context, req *pb.KvsDeleteRequest) (*pb.KvsDeleteResponse, error) {
	return s.parent.storageManager.KvsDelete(c, req)
}

func (s *clusterServer) KvsGetKeys(c context.Context, req *pb.KvsAllKeysRequest) (*pb.KvsAllKeysResponse, error) {
	return s.parent.storageManager.KvsGetKeys(c, req)
}

func (s *clusterServer) TSSave(c context.Context, req *pb.TSStoreRequest) (*pb.TSStoreResponse, error) {
	return s.parent.storageManager.TSSave(c, req)
}

func (s *clusterServer) TSRetrieve(c context.Context, req *pb.TSRetrieveRequest) (*pb.TSRetrieveResponse, error) {
	return s.parent.storageManager.TSRetrieve(c, req)
}

func (s *clusterServer) TSAvailability(c context.Context, req *pb.TSAvailabilityRequest) (*pb.TSAvailabilityResponse, error) {
	return s.parent.storageManager.TSAvailability(c, req)
}