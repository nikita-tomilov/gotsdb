package tss

import (
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"sync"
	"time"
)

type InMemTSS struct {
	isRunning          bool
	data               map[string]TSforDatasource
	lock               sync.Mutex
	periodBetweenWipes time.Duration
	maxEntriesPerTag   int
}

func (f *InMemTSS) InitStorage() {
	f.isRunning = true
	if f.periodBetweenWipes == 0*time.Second {
		f.periodBetweenWipes = time.Second * 5
	}
	if f.maxEntriesPerTag == 0 {
		f.maxEntriesPerTag = 21600
	}
	f.data = make(map[string]TSforDatasource)
}

func (f *InMemTSS) CloseStorage() {
	f.isRunning = false
}

func (f *InMemTSS) Save(dataSource string, data map[string]*pb.TSPoints, expirationMillis uint64) {
	f.lock.Lock()
	if !f.contains(dataSource) {
		dataForDataSource := TSforDatasource{PeriodBetweenWipes: f.periodBetweenWipes, MaxEntriesPerTag: f.maxEntriesPerTag}
		dataForDataSource.Init()
		f.data[dataSource] = dataForDataSource
	}
	f.dataForDataSource(dataSource).SaveData(data, expirationMillis)
	f.lock.Unlock()
}

func (f *InMemTSS) SaveBatch(dataSource string, data []*pb.TSPoint, expirationMillis uint64) {
	f.lock.Lock()
	if !f.contains(dataSource) {
		dataForDataSource := TSforDatasource{PeriodBetweenWipes: f.periodBetweenWipes, MaxEntriesPerTag: f.maxEntriesPerTag}
		dataForDataSource.Init()
		f.data[dataSource] = dataForDataSource
	}
	f.dataForDataSource(dataSource).SaveDataBatch(data, expirationMillis)
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

func (f *InMemTSS) dataForDataSource(dataSource string) *TSforDatasource {
	dataForDataSource := f.data[dataSource]
	return &dataForDataSource
}
