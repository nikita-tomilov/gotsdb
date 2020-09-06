package tss

import (
	"fmt"
	"github.com/nikita-tomilov/golsm"
	"github.com/nikita-tomilov/golsm/dto"
	"github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/utils"
	"sync"
	"time"
)

type StoragesForDatasource struct {
	storageReader *golsm.StorageReader
	storageWriter *golsm.StorageWriter
}

type LSMTSS struct {
	Path                        string `summer.property:"ts.filePath|/tmp/gotsdb/tss"`
	CommitlogFlushPeriodSeconds int    `summer.property:"tslsm.commitlogFlushPeriodSeconds|10"`
	MemtExpirationPeriodSeconds int    `summer.property:"tslsm.memtExpirationPeriodSeconds|10"`
	forDatasource               map[string]StoragesForDatasource
	mutex                       *sync.Mutex
}

func (lsm *LSMTSS) getOrInitStorage(datasource string) StoragesForDatasource {
	val, ok := lsm.forDatasource[datasource]
	if ok {
		return val
	}

	lsm.mutex.Lock()
	val, ok = lsm.forDatasource[datasource]
	if ok {
		lsm.mutex.Unlock()
		return val
	}

	rootPath := lsm.Path + "/" + datasource
	storageReader, storageWriter := golsm.InitStorage(
		rootPath+"/commitlog",
		10,
		time.Duration(lsm.CommitlogFlushPeriodSeconds)*time.Second,
		time.Duration(lsm.MemtExpirationPeriodSeconds)*time.Second,
		rootPath+"/sst",
		100)
	lsm.forDatasource[datasource] = StoragesForDatasource{storageReader: storageReader, storageWriter: storageWriter}

	lsm.mutex.Unlock()
	return lsm.forDatasource[datasource]
}

func (lsm *LSMTSS) InitStorage() {
	lsm.mutex = &sync.Mutex{}
	lsm.forDatasource = make(map[string]StoragesForDatasource)
}

func (lsm *LSMTSS) CloseStorage() {
	//nothing
}

func (lsm *LSMTSS) Save(dataSource string, data map[string]*proto.TSPoints, expirationMillis uint64) {
	converted := make(map[string][]dto.Measurement)
	expireAt := utils.GetNowMillis() + expirationMillis
	if expirationMillis == 0 {
		expireAt = 0
	}
	for k, v := range data {
		converted[k] = convertTSPtoMeasurement(v)
	}
	lsm.getOrInitStorage(dataSource).storageWriter.Store(converted, expireAt)
}

func (lsm *LSMTSS) Retrieve(dataSource string, tags []string, fromTimestamp uint64, toTimestamp uint64) map[string]*proto.TSPoints {
	retrieved := lsm.getOrInitStorage(dataSource).storageReader.Retrieve(tags, fromTimestamp, toTimestamp)
	deconverted := make(map[string]*proto.TSPoints)
	for tag, values := range retrieved {
		deconverted[tag] = convertMeasurementToTSP(values)
	}
	return deconverted
}

func convertTSPtoMeasurement(points *proto.TSPoints) []dto.Measurement {
	ans := make([]dto.Measurement, len(points.Points))
	i := 0
	for ts, val := range points.Points {
		ans[i] = dto.Measurement{Timestamp: ts, Value: utils.Float64ToByte(val)}
		i++
	}
	return ans
}

func convertMeasurementToTSP(measurements []dto.Measurement) *proto.TSPoints {
	ans := make(map[uint64]float64)
	for _, m := range measurements {
		ans[m.Timestamp] = utils.ByteToFloat64(m.Value)
	}
	return &proto.TSPoints{Points: ans}
}

func (lsm *LSMTSS) Availability(dataSource string, fromTimestamp uint64, toTimestamp uint64) []*proto.TSAvailabilityChunk {
	from, to := lsm.getOrInitStorage(dataSource).storageReader.Availability()
	if (from == 0) && (to == 0) {
		return []*proto.TSAvailabilityChunk{}
	}
	return []*proto.TSAvailabilityChunk{{FromTimestamp: from, ToTimestamp: to}}
}

func (lsm *LSMTSS) String() string {
	return fmt.Sprintf("LSM-based storage over the root dir %s", lsm.Path)
}
