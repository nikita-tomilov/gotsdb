package tss

import (
	"fmt"
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/utils"
	"strings"
	"time"
)

type MeasurementMeta struct {
	Id  uint `gorm:"primaryKey"`
	Tag string
}

type MeasurementDataSource struct {
	Id         uint `gorm:"primaryKey"`
	DataSource string
}

type Measurement struct {
	DataSourceKey uint   `gorm:"index"`
	MetaKey       uint   `gorm:"index"`
	Ts            uint64 `gorm:"index"`
	Value         float64
	ExpireAt      uint64 `gorm:"index"`
}

type SqlWrapper interface {
	InitDatabase()
	Execute(query string)
	DeleteOnExpiration()
	CreateMetaKey(tag string)
	GetMetaKey(tag string) (uint, error)
	CreateDataSourceKey(dataSource string)
	GetDataSourceKey(dataSource string) (uint, error)
	GetTwoTimestamps(query string) (uint64, uint64)
	GetMeasurementsForTag(query string) map[uint64]float64
	Close()
}

type AbstractSQLTSS struct {
	sqlWrapper         SqlWrapper
	periodBetweenWipes time.Duration
	isRunning          bool
	dataSourceCache    map[string]uint
	metaCache          map[string]uint
}

func (sq *AbstractSQLTSS) Init() {
	sq.sqlWrapper.InitDatabase()
	sq.metaCache = make(map[string]uint)
	sq.dataSourceCache = make(map[string]uint)
	go func(sq *AbstractSQLTSS) {
		time.Sleep(sq.periodBetweenWipes)
		for sq.isRunning {
			sq.expirationCycle()
			time.Sleep(sq.periodBetweenWipes)
		}
	}(sq)
}

func (sq *AbstractSQLTSS) Close() {
	sq.sqlWrapper.Close()
	sq.isRunning = false
}

func (sq *AbstractSQLTSS) getMetaKey(tag string) uint {
	cached, existsInCache := sq.metaCache[tag]
	if existsInCache {
		return cached
	}
	key, err := sq.sqlWrapper.GetMetaKey(tag)
	if err != nil {
		sq.sqlWrapper.CreateMetaKey(tag)
		return sq.getMetaKey(tag)
	}
	sq.metaCache[tag] = key
	return key
}

func (sq *AbstractSQLTSS) getDataSourceKey(ds string) uint {
	cached, existsInCache := sq.dataSourceCache[ds]
	if existsInCache {
		return cached
	}
	key, err := sq.sqlWrapper.GetDataSourceKey(ds)
	if err != nil {
		sq.sqlWrapper.CreateDataSourceKey(ds)
		return sq.getDataSourceKey(ds)
	}
	sq.dataSourceCache[ds] = key
	return key
}

func (sq *AbstractSQLTSS) internalSaveBatch(batch []Measurement, actualLen int) {
	sb := strings.Builder{}
	sb.WriteString("BEGIN TRANSACTION;")
	i := 0
	for i < actualLen {
		entry := batch[i]
		sb.WriteString(fmt.Sprintf("INSERT INTO measurements VALUES(%d, %d, %d, %f, %d);", entry.DataSourceKey, entry.MetaKey, entry.Ts, entry.Value, entry.ExpireAt))
		i++
	}
	sb.WriteString("COMMIT;")
	rq := sb.String()
	sq.sqlWrapper.Execute(rq)
}

func (sq *AbstractSQLTSS) SaveBatch(dataSource string, data []*pb.TSPoint, expirationMillis uint64) {
	now := utils.GetNowMillis()
	expireAt := now + expirationMillis
	if expirationMillis == 0 {
		expireAt = 0
	}
	const batchLength = 100
	i := 0
	measurements := make([]Measurement, batchLength)

	for _, point := range data {
		tag := point.Tag
		ts := point.Timestamp
		value := point.Value
		metaKey := sq.getMetaKey(tag)
		dsKey := sq.getDataSourceKey(dataSource)
		meas := Measurement{DataSourceKey: dsKey, MetaKey: metaKey, Ts: ts, Value: value, ExpireAt: expireAt}
		measurements[i] = meas
		i += 1
		if i == batchLength {
			sq.internalSaveBatch(measurements, i)
			i = 0
		}
	}
	if i > 0 {
		sq.internalSaveBatch(measurements, i)
	}
}

func (sq *AbstractSQLTSS) Save(dataSource string, data map[string]*pb.TSPoints, expirationMillis uint64) {
	now := utils.GetNowMillis()
	expireAt := now + expirationMillis
	if expirationMillis == 0 {
		expireAt = 0
	}
	const batchLength = 100
	i := 0
	measurements := make([]Measurement, batchLength)
	for tag, values := range data {
		metaKey := sq.getMetaKey(tag)
		dsKey := sq.getDataSourceKey(dataSource)
		for ts, value := range values.Points {
			meas := Measurement{DataSourceKey: dsKey, MetaKey: metaKey, Ts: ts, Value: value, ExpireAt: expireAt}
			measurements[i] = meas
			i += 1
			if i == batchLength {
				sq.internalSaveBatch(measurements, i)
				i = 0
			}
		}
	}
	if i > 0 {
		sq.internalSaveBatch(measurements, i)
	}
}

func (sq *AbstractSQLTSS) Retrieve(dataSource string, tags []string, fromTimestamp uint64, toTimestamp uint64) map[string]*pb.TSPoints {
	ans := make(map[string]*pb.TSPoints)
	for _, tag := range tags {
		rq := fmt.Sprintf("SELECT ts, value, expire_at FROM measurements WHERE data_source_key = %d AND meta_key = %d AND ts >= %d AND ts <= %d", sq.getDataSourceKey(dataSource), sq.getMetaKey(tag), fromTimestamp, toTimestamp)
		ansForTag := sq.sqlWrapper.GetMeasurementsForTag(rq)
		ans[tag] = &pb.TSPoints{Points: ansForTag}
	}
	return ans
}

func (sq *AbstractSQLTSS) Availability(dataSource string, fromTimestamp uint64, toTimestamp uint64) []*pb.TSAvailabilityChunk {
	now := utils.GetNowMillis()
	rq := fmt.Sprintf("SELECT min(ts), max(ts) FROM measurements WHERE data_source_key = %d AND (expire_at == 0 OR expire_at > %d);", sq.getDataSourceKey(dataSource), now)
	min, max := sq.sqlWrapper.GetTwoTimestamps(rq)
	if min >= max {
		ans := make([]*pb.TSAvailabilityChunk, 0)
		return ans
	}
	min = utils.Max(fromTimestamp, min)
	max = utils.Min(toTimestamp, max)
	ans := make([]*pb.TSAvailabilityChunk, 1)
	ans[0] = &pb.TSAvailabilityChunk{FromTimestamp: min, ToTimestamp: max}
	return ans
}

func (sq *AbstractSQLTSS) expirationCycle() {
	sq.sqlWrapper.DeleteOnExpiration()
}
