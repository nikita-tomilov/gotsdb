package cluster

import (
	"context"
	"github.com/google/uuid"
	log "github.com/jeanphorn/log4go"
	pb "github.com/programmer74/gotsdb/proto"
	"strings"
	"time"
)

type Manager struct {
	GrpcClusterServer    *interface{} `summer:"*cluster.GrpcClusterServer"`
	KnownNodesParam      string       `summer.property:"cluster.knownNodes|localhost:5123"`
	IncomingConnections  map[string]pb.Node
	OutcomingConnections map[string]*GrpcClusterClient
	MyUUID               string
	MyHostPort           string `summer.property:"cluster.listenAddress|:5300"`
}

func (m *Manager) getGrpcServer() *GrpcClusterServer {
	s := m.GrpcClusterServer
	s2 := (*s).(*GrpcClusterServer)
	return s2
}

func (m *Manager) StartClustering() {
	m.getGrpcServer().manager = m
	m.getGrpcServer().Start()

	m.MyUUID = uuid.New().String()

	m.IncomingConnections = make(map[string]pb.Node)
	m.OutcomingConnections = make(map[string]*GrpcClusterClient)

	urlsToProbe := strings.Split(m.KnownNodesParam, ";")
	for _, url := range urlsToProbe {
		log.Warn("Should probe %s", url)
		client := GrpcClusterClient{myUUID: m.MyUUID, myHostPort: m.MyHostPort, targetHostPort: url, nodeDisconnectedCallback: m.DeleteKnownNode}
		client.BeginProbing()
		m.OutcomingConnections[url] = &client
	}

	m.startPrintingKnownNodes()
}

func (m *Manager) startPrintingKnownNodes() {
	go func() {
		for true {
			time.Sleep(5 * time.Second)
			log.Info("Known nodes: %d", len(m.IncomingConnections))
			for _, v := range m.IncomingConnections {
				log.Info("- %s at %s", v.Uuid, v.ConnectionString)
			}

			log.Info("Consensus testing: ")
			lenOfKnownNodesOnThisMachine := len(m.IncomingConnections)
			for k, v := range m.OutcomingConnections {
				if v.IsConnected() {
					ans, _ := v.GetGrpcChannel().GetAliveNodes(context.TODO(), &pb.Void{})
					lenOfKnownNodesOnAnotherMachine := len(ans.AliveNodes)
					if lenOfKnownNodesOnAnotherMachine == lenOfKnownNodesOnThisMachine {
						log.Info(" - %s: Consensus reached", k)
					} else {
						log.Warn(" - %s: Consensus NOT reached: i think %d nodes, other thinks %d nodes", k, lenOfKnownNodesOnThisMachine, lenOfKnownNodesOnAnotherMachine)
					}
				} else {
					log.Info(" - %s: Connection problems", k)
				}
			}

			log.Info("Known connections to other nodes: %d", len(m.GetKnownOutboundConnections()))
		}
	}()
}

func (m *Manager) AddKnownNode(n *pb.Node) {
	m.IncomingConnections[n.ConnectionString] = *n
}

func (m *Manager) GetKnownNodes() []*pb.Node {
	var arr []*pb.Node
	for _, v := range m.IncomingConnections {
		arr = append(arr, &v)
	}
	return arr
}

func (m *Manager) DeleteKnownNode(targetHostPort string) {
	delete(m.IncomingConnections, targetHostPort)
}

func (m *Manager) GetKnownOutboundConnections() []*GrpcClusterClient {
	var arr []*GrpcClusterClient
	for k, _ := range m.IncomingConnections {
		if k != m.MyHostPort {
			arr = append(arr, m.OutcomingConnections[k])
		}
	}
	return arr
}