package cluster

import (
	"context"
	log "github.com/jeanphorn/log4go"
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"google.golang.org/grpc"
	"sync/atomic"
	"time"
)

type NodeDisconnectedCallback func(string)

type GrpcClusterClient struct {
	isConnected              atomic.Value
	myUUID                   string
	myHostPort               string
	targetHostPort           string
	nodeDisconnectedCallback NodeDisconnectedCallback
	grpcChannel              atomic.Value
}

func (c *GrpcClusterClient) IsConnected() bool {
	return c.isConnected.Load().(bool)
}

func (c *GrpcClusterClient) GetGrpcChannel() pb.ClusterClient {
	link := c.grpcChannel.Load().(pb.ClusterClient)
	return link
}

func (c *GrpcClusterClient) BeginProbing() {
	go func() {
		c.isConnected.Store(false)

		for true {
			opts := []grpc.DialOption{
				grpc.WithInsecure(),
				grpc.WithBlock(),
			}
			conn, err := grpc.Dial(c.targetHostPort, opts...)

			if err != nil {
				time.Sleep(1 * time.Second)
				log.Warn("Failed to dial %s: %v", c.targetHostPort, err)
				continue
			}

			client := pb.NewClusterClient(conn)
			client.Hello(context.TODO(), &pb.HelloRequest{Iam: &pb.Node{ConnectionString: c.myHostPort, Uuid: c.myUUID}})
			isConnected := true

			c.isConnected.Store(true)
			c.grpcChannel.Store(client)

			for isConnected {
				time.Sleep(1 * time.Second)
				pingRq := pb.PingRequest{Payload: make([]byte, 0)}
				_, err := client.Ping(context.TODO(), &pingRq)
				if err != nil {
					log.Warn("Connection broken to %s: %v", c.targetHostPort, err)
					isConnected = false

					c.isConnected.Store(false)

					conn.Close()
					c.nodeDisconnectedCallback(c.targetHostPort)
				}
				log.Debug("Outer connection to %s is alive", c.targetHostPort)
			}
		}
	}()
}
