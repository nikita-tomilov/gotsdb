package tss

import (
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/services/storage/tss/dummytss"
	"sync"
	"time"
)

type InMemTSS struct {
	data map[string]dummytss.TSforDatasource
	lock sync.Mutex
}

func (f *InMemTSS) InitStorage() {
	go func() {
		for true {
			f.lock.Lock()
			for _, d := range f.data {
				d.ExpirationCycle()
			}
			f.lock.Unlock()
			time.Sleep(time.Second * 5)
		}
	}()
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

func (f *InMemTSS) contains(dataSource string) bool {
	_, found := f.data[dataSource]
	return found
}

func (f *InMemTSS) dataForDataSource(dataSource string) *dummytss.TSforDatasource {
	dataForDataSource := f.data[dataSource]
	return &dataForDataSource
}