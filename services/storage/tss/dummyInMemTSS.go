package tss

import (
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/services/storage/tss/dummytss"
	"sync"
	"time"
)

type InMemTSS struct {
	isRunning bool
	data map[string]dummytss.TSforDatasource
	lock sync.Mutex
	periodBetweenWipes time.Duration
}

func (f *InMemTSS) InitStorage() {
	f.isRunning = true
	if f.periodBetweenWipes == 0 * time.Second {
		f.periodBetweenWipes = time.Second * 5
	}
	f.data = make(map[string]dummytss.TSforDatasource)
	go func(s *InMemTSS) {
		for s.isRunning {
			s.lock.Lock()
			for _, d := range s.data {
				d.ExpirationCycle()
			}
			s.lock.Unlock()
			time.Sleep(s.periodBetweenWipes)
		}
	}(f)
}

func (f *InMemTSS) CloseStorage() {
	f.isRunning = false
}

func (f *InMemTSS) Save(dataSource string, data map[string]*pb.TSPoints, expirationMillis uint64) {
	f.lock.Lock()
	if !f.contains(dataSource) {
		dataForDataSource := dummytss.TSforDatasource{}
		dataForDataSource.Init()
		f.data[dataSource] = dataForDataSource
	}
	f.dataForDataSource(dataSource).SaveData(data, expirationMillis)
	f.lock.Unlock()
}

func (f *InMemTSS) Retrieve(dataSource string, tags []string, fromTimestamp uint64, toTimestamp uint64) map[string]*pb.TSPoints {
	f.lock.Lock()
	ans := make(map[string]*pb.TSPoints)
	if f.contains(dataSource) {
		ans = f.dataForDataSource(dataSource).GetData(tags, fromTimestamp, toTimestamp)
	}
	f.lock.Unlock()
	return ans
}

func (f *InMemTSS) Availability(dataSource string, fromTimestamp uint64, toTimestamp uint64) []*pb.TSAvailabilityChunk {
	f.lock.Lock()
	var ans []*pb.TSAvailabilityChunk
	if f.contains(dataSource) {
		ans = f.dataForDataSource(dataSource).Availability(fromTimestamp, toTimestamp)
	} else {
		ans = make([]*pb.TSAvailabilityChunk, 0)
	}
	f.lock.Unlock()
	return ans
}

func (f *InMemTSS) String() string {
	return "Dummy In-Memory map-based TSS"
}

func (f *InMemTSS) contains(dataSource string) bool {
	_, found := f.data[dataSource]
	return found
}

func (f *InMemTSS) dataForDataSource(dataSource string) *dummytss.TSforDatasource {
	dataForDataSource := f.data[dataSource]
	return &dataForDataSource
}