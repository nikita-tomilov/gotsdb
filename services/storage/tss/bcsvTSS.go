package tss

import (
	"github.com/nikita-tomilov/gotsdb/proto"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type BCSVTSS struct {
	Path               string `summer.property:"tss.filePath|/tmp/gotsdb/tss"`
	data               map[string]BCSVTSforDatasource
	periodBetweenWipes time.Duration
	lock               sync.Mutex
	isRunning          bool
}

func (c *BCSVTSS) InitStorage() {
	os.MkdirAll(c.Path, os.ModePerm)
	c.isRunning = true
	if c.periodBetweenWipes == 0*time.Second {
		c.periodBetweenWipes = time.Second * 5
	}
	c.data = make(map[string]BCSVTSforDatasource)
	datasources := c.getAllDirectories()
	for _, ds := range datasources {
		c.initDataSource(ds)
	}
	go func(c *BCSVTSS) {
		time.Sleep(c.periodBetweenWipes)
		for c.isRunning {
			c.expirationCycle()
			time.Sleep(c.periodBetweenWipes)
		}
	}(c)
}

func (c *BCSVTSS) CloseStorage() {
	c.isRunning = false
}

func (c *BCSVTSS) Save(dataSource string, data map[string]*proto.TSPoints, expirationMillis uint64) {
	c.lock.Lock()
	if !c.contains(dataSource) {
		c.initDataSource(dataSource)
	}
	c.dataForDataSource(dataSource).SaveData(data, expirationMillis)
	c.lock.Unlock()
}

func (c *BCSVTSS) SaveBatch(dataSource string, data []*proto.TSPoint, expirationMillis uint64) {
	c.lock.Lock()
	if !c.contains(dataSource) {
		c.initDataSource(dataSource)
	}
	c.dataForDataSource(dataSource).SaveDataBatch(data, expirationMillis)
	c.lock.Unlock()
}

func (c *BCSVTSS) Retrieve(dataSource string, tags []string, fromTimestamp uint64, toTimestamp uint64) map[string]*proto.TSPoints {
	c.lock.Lock()
	ans := make(map[string]*proto.TSPoints)
	if c.contains(dataSource) {
		ans = c.dataForDataSource(dataSource).GetData(tags, fromTimestamp, toTimestamp)
	}
	c.lock.Unlock()
	return ans
}

func (c *BCSVTSS) Availability(dataSource string, fromTimestamp uint64, toTimestamp uint64) []*proto.TSAvailabilityChunk {
	c.lock.Lock()
	var ans []*proto.TSAvailabilityChunk
	if c.contains(dataSource) {
		ans = c.dataForDataSource(dataSource).Availability(fromTimestamp, toTimestamp)
	} else {
		ans = make([]*proto.TSAvailabilityChunk, 0)
	}
	c.lock.Unlock()
	return ans
}

func (c *BCSVTSS) String() string {
	return "Binary CSV-based disk TSS at " + c.Path
}

func (c *BCSVTSS) contains(dataSource string) bool {
	_, found := c.data[dataSource]
	return found
}

func (c *BCSVTSS) dataForDataSource(dataSource string) *BCSVTSforDatasource {
	dataForDataSource := c.data[dataSource]
	return &dataForDataSource
}

func (c *BCSVTSS) initDataSource(dataSource string) {
	dataForDataSource := BCSVTSforDatasource{DatasourcePath: c.Path + "/" + dataSource}
	dataForDataSource.Init()
	c.data[dataSource] = dataForDataSource
}

func (c *BCSVTSS) expirationCycle() {
	c.lock.Lock()
	for _, fords := range c.data {
		fords.ExpirationCycle()
	}
	c.lock.Unlock()
}

func (c *BCSVTSS) getAllDirectories() []string {
	var ans []string
	files, err := ioutil.ReadDir(c.Path)
	if err != nil {
		panic(err)
	}
	for _, path := range files {
		if path.IsDir() {
			ans = append(ans, path.Name())
		}
	}
	return ans
}