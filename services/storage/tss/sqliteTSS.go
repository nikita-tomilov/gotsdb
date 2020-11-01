package tss

import (
	"database/sql"
	"fmt"
	log "github.com/jeanphorn/log4go"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nikita-tomilov/gotsdb/proto"
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"math"
	"os"
	"strings"
	"time"
)

type SqliteTSS struct {
	db                 *gorm.DB
	Path               string `summer.property:"tss.filePath|/tmp/gotsdb/tss"`
	dbFilePath         string
	periodBetweenWipes time.Duration
	isRunning          bool
}

type Measurement struct {
	DataSource string `gorm:"index"`
	Tag        string `gorm:"index"`
	Ts         uint64 `gorm:"index"`
	Value      float64
	ExpireAt   uint64
}

func (sq *SqliteTSS) InitStorage() {
	_ = os.MkdirAll(sq.Path, os.ModePerm)
	sq.dbFilePath = sq.Path + "/db.bin"
	db, err := gorm.Open(sqlite.Open(sq.dbFilePath), &gorm.Config{})
	if err != nil {
		panic("Unable to instantiate db " + err.Error())
	}
	sq.db = db
	sq.createTableIfNotExists()

	sq.isRunning = true
	if sq.periodBetweenWipes == 0*time.Second {
		sq.periodBetweenWipes = time.Second * 5
	}
	go func(sq *SqliteTSS) {
		time.Sleep(sq.periodBetweenWipes)
		for sq.isRunning {
			sq.expirationCycle()
			time.Sleep(sq.periodBetweenWipes)
		}
	}(sq)
}

func (sq *SqliteTSS) createTableIfNotExists() {
	err := sq.db.AutoMigrate(&Measurement{})
	if err != nil {
		log.Error("Error in DB: " + err.Error())
		panic(err)
	}
}

func (sq *SqliteTSS) CloseStorage() {
	sqlDb, err := sq.db.DB()
	if err != nil {
		log.Error("Error in DB Close: " + err.Error())
	} else {
		sqlDb.Close()
	}
	sq.isRunning = false
}

func (sq *SqliteTSS) saveBatch(batch []Measurement) {
	sb := strings.Builder{}
	sb.WriteString("BEGIN TRANSACTION;")
	for _, entry := range batch {
		sb.WriteString(fmt.Sprintf("INSERT INTO measurements VALUES(\"%s\", \"%s\", %d, %f, %d);", entry.DataSource, entry.Tag, entry.Ts, entry.Value, entry.ExpireAt))
	}
	sb.WriteString("COMMIT;")
	rq := sb.String()
	db, _ := sq.db.DB()
	_, err := db.Exec(rq)
	if err != nil {
		log.Error("Error in DB in SaveBatch: " + err.Error())
	}
}

func (sq *SqliteTSS) toBatches(total []Measurement, batchSize int) [][]Measurement {
	i := 0
	ans := make([][]Measurement, 0)
	for i < len(total) {
		j := utils.MinInt(i + batchSize, len(total))
		ans = append(ans, total[i:j])
		i = i + batchSize
	}
	return ans
}

func (sq *SqliteTSS) Save(dataSource string, data map[string]*proto.TSPoints, expirationMillis uint64) {
	now := utils.GetNowMillis()
	expireAt := now + expirationMillis
	measurements := make([]Measurement, 0)
	for tag, values := range data {
		measurementsForTag := make([]Measurement, len(values.Points))
		i := 0
		for ts, value := range values.Points {
			meas := Measurement{DataSource: dataSource, Tag: tag, Ts: ts, Value: value, ExpireAt: expireAt}
			measurementsForTag[i] = meas
			i += 1
		}
		measurements = append(measurements, measurementsForTag...)
	}
	for _, batch := range sq.toBatches(measurements, 100) {
		sq.saveBatch(batch)
	}
}

func (sq *SqliteTSS) Retrieve(dataSource string, tags []string, fromTimestamp uint64, toTimestamp uint64) map[string]*proto.TSPoints {
	ans := make(map[string]*proto.TSPoints)

	for _, tag := range tags {
		ansForTag := make(map[uint64]float64)
		rq := fmt.Sprintf("SELECT ts, value FROM measurements WHERE data_source = \"%s\" AND tag = \"%s\" AND ts >= %d AND ts <= %d", dataSource, tag, fromTimestamp, toTimestamp)

		db, _ := sq.db.DB()
		rows, err := db.Query(rq)
		if err != nil {
			log.Error("Error in DB in Retrieve: " + err.Error())
		} else {
			ts := uint64(0)
			val := 0.0
			for rows.Next() {
				err := rows.Scan(&ts, &val)
				if err != nil {
					log.Error(err)
				} else {
					ansForTag[ts] = val
				}
			}
		}

		ans[tag] = &pb.TSPoints{Points: ansForTag}
	}

	return ans
}

func (sq *SqliteTSS) Availability(dataSource string, fromTimestamp uint64, toTimestamp uint64) []*proto.TSAvailabilityChunk {
	rq := fmt.Sprintf("SELECT min(ts), max(ts) FROM measurements WHERE data_source = \"%s\";", dataSource)
	res, err := sq.db.Raw(rq).Rows()
	min := uint64(math.MaxUint64)
	max := uint64(0)
	if err != nil {
		log.Error("Error in DB in Availability: " + err.Error())
		ans := make([]*proto.TSAvailabilityChunk, 0)
		return ans
	} else {
		for res.Next() {
			min2 := sql.NullInt64{}
			max2 := sql.NullInt64{}
			err := res.Scan(&min2, &max2)
			if err != nil {
				log.Error(err)
			} else {
				if min2.Valid && max2.Valid {
					min = utils.Min(min, uint64(min2.Int64))
					max = utils.Max(max, uint64(max2.Int64))
				}
			}
		}
	}
	if min >= max {
		ans := make([]*proto.TSAvailabilityChunk, 0)
		return ans
	}
	min = utils.Max(fromTimestamp, min)
	max = utils.Min(toTimestamp, max)
	ans := make([]*proto.TSAvailabilityChunk, 1)
	ans[0] = &proto.TSAvailabilityChunk{FromTimestamp: min, ToTimestamp: max}
	return ans
}

func (sq *SqliteTSS) String() string {
	return fmt.Sprintf("SqliteTSS on file %s", sq.dbFilePath)
}

func (sq *SqliteTSS) expirationCycle() {
	now := utils.GetNowMillis()
	err := sq.db.Begin().Delete(Measurement{}, "expire_at <= ?", now).Commit().Error

	if err != nil {
		log.Error("Error in DB in expirationCycle: " + err.Error())
		sq.isRunning = false
	}
}
