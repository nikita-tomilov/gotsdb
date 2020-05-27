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
}

func (s *GrpcClusterServer) Start() {
	log.Warn("Starting to listen for other nodes at '%s'", s.ListenAddress)
	listener, err := net.Listen("tcp", s.ListenAddress)

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterClusterServer(grpcServer, &clusterServer{})
	grpcServer.Serve(listener)
}

//TODO: pass errors from storage level to grpc level?

type clusterServer struct{

}

func (s clusterServer) Ping(c context.Context, rq *pb.PingRequest) (*pb.PingResponse, error) {
	panic("implement me")
}

func (s clusterServer) GetAliveNodes(c context.Context, v *pb.Void) (*pb.AliveNodesResponse, error) {
	panic("implement me")
}

