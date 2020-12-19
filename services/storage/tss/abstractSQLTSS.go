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
	Key string
}

type Measurement struct {
	DataSource string
	Key        uint   `gorm:"index"`
	Ts         uint64 `gorm:"index"`
	Value      float64
	ExpireAt   uint64
}

type SqlWrapper interface {
	InitDatabase()
	Execute(query string)
	DeleteOnExpiration()
	CreateMetaKey(tag string)
	GetMetaKey(tag string) (uint, error)
	GetTwoTimestamps(query string) (uint64, uint64)
	GetMeasurementsForTag(query string) map[uint64]float64
	Close()
}

type AbstractSQLTSS struct {
	sqlWrapper         SqlWrapper
	periodBetweenWipes time.Duration
	isRunning          bool
}

func (sq *AbstractSQLTSS) Init() {
	sq.sqlWrapper.InitDatabase()
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

func (sq *AbstractSQLTSS) getKeyByDsAndTag(ds string, tag string) uint {
	k := ds + "_" + tag
	key, err := sq.sqlWrapper.GetMetaKey(k)
	if err != nil {
		sq.sqlWrapper.CreateMetaKey(k)
		return sq.getKeyByDsAndTag(ds, tag)
	}
	return key
}

func (sq *AbstractSQLTSS) saveBatch(batch []Measurement, actualLen int) {
	sb := strings.Builder{}
	sb.WriteString("BEGIN TRANSACTION;")
	i := 0
	for i < actualLen {
		entry := batch[i]
		sb.WriteString(fmt.Sprintf("INSERT INTO measurements VALUES(\"%s\", %d, %d, %f, %d);", entry.DataSource, entry.Key, entry.Ts, entry.Value, entry.ExpireAt))
		i++
	}
	sb.WriteString("COMMIT;")
	rq := sb.String()
	sq.sqlWrapper.Execute(rq)
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
		key := sq.getKeyByDsAndTag(dataSource, tag)
		for ts, value := range values.Points {
			meas := Measurement{DataSource: dataSource, Key: key, Ts: ts, Value: value, ExpireAt: expireAt}
			measurements[i] = meas
			i += 1
			if i == batchLength {
				sq.saveBatch(measurements, i)
				i = 0
			}
		}
	}
	if i > 0 {
		sq.saveBatch(measurements, i)
	}
}

func (sq *AbstractSQLTSS) Retrieve(dataSource string, tags []string, fromTimestamp uint64, toTimestamp uint64) map[string]*pb.TSPoints {
	ans := make(map[string]*pb.TSPoints)
	for _, tag := range tags {
		rq := fmt.Sprintf("SELECT ts, value, expire_at FROM measurements WHERE meta_key = %d AND ts >= %d AND ts <= %d", sq.getKeyByDsAndTag(dataSource, tag), fromTimestamp, toTimestamp)
		ansForTag := sq.sqlWrapper.GetMeasurementsForTag(rq)
		ans[tag] = &pb.TSPoints{Points: ansForTag}
	}
	return ans
}

func (sq *AbstractSQLTSS) Availability(dataSource string, fromTimestamp uint64, toTimestamp uint64) []*pb.TSAvailabilityChunk {
	now := utils.GetNowMillis()
	rq := fmt.Sprintf("SELECT min(ts), max(ts) FROM measurements WHERE data_source IN (\"%s\") AND (expire_at == 0 OR expire_at > %d);", dataSource, now)
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
