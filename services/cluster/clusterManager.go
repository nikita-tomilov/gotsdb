package cluster

type Manager struct {
	GrpcClusterServer *interface{} `summer:"*cluster.GrpcClusterServer"`
}

func (m *Manager) getGrpcServer() *GrpcClusterServer {
	s := m.GrpcClusterServer
	s2 := (*s).(*GrpcClusterServer)
	return s2
}

func (m Manager) StartClustering() {
	go m.getGrpcServer().Start()
}