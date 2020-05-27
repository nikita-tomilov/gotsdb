package cluster

import (
	log "github.com/jeanphorn/log4go"
	pb "github.com/programmer74/gotsdb/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
)

type GrpcClusterServer struct{
	ListenAddress string `summer.property:"cluster.listenAddress|:5300"`
	manager *Manager
}

func (s *GrpcClusterServer) Start() {
	log.Warn("Starting to listen for other nodes at '%s'", s.ListenAddress)
	listener, err := net.Listen("tcp", s.ListenAddress)

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterClusterServer(grpcServer, &clusterServer{parent:s})
	go grpcServer.Serve(listener)
}

//TODO: pass errors from storage level to grpc level?

type clusterServer struct{
	parent *GrpcClusterServer
}

func (s clusterServer) Hello(c context.Context, rq *pb.HelloRequest) (*pb.AliveNodesResponse, error) {
	s.parent.manager.AddKnownNode(rq.Iam)
	return &pb.AliveNodesResponse{AliveNodes:s.parent.manager.GetKnownNodes()}, nil
}

func (s clusterServer) Ping(c context.Context, rq *pb.PingRequest) (*pb.PingResponse, error) {
	log.Debug("Got ping rq")
	return &pb.PingResponse{Payload:rq.Payload}, nil
}

func (s clusterServer) GetAliveNodes(c context.Context, v *pb.Void) (*pb.AliveNodesResponse, error) {
	return &pb.AliveNodesResponse{AliveNodes:s.parent.manager.GetKnownNodes()}, nil
}

