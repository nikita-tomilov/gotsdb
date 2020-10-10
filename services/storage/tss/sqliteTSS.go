package tss

import (
	"database/sql"
	"fmt"
	log "github.com/jeanphorn/log4go"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nikita-tomilov/gotsdb/proto"
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/utils"
	"math"
	"os"
	"strings"
	"time"
)

type SqliteTSS struct {
	db                 *sql.DB
	Path               string `summer.property:"tss.filePath|/tmp/gotsdb/tss"`
	dbFilePath         string
	periodBetweenWipes time.Duration
	isRunning          bool
}

func (sq *SqliteTSS) InitStorage() {
	_ = os.MkdirAll(sq.Path, os.ModePerm)
	sq.dbFilePath = sq.Path + "/db.bin"
	db, err := sql.Open("sqlite3", sq.dbFilePath)
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
		for sq.isRunning {
			sq.expirationCycle()
			time.Sleep(sq.periodBetweenWipes)
		}
	}(sq)
}

func (sq *SqliteTSS) createTableIfNotExists() {
	_, err := sq.db.Exec(`
		BEGIN TRANSACTION;
			CREATE TABLE IF NOT EXISTS RawData (ds string, tag string, ts uint64, value float64, expat uint64);
			CREATE INDEX IF NOT EXISTS RawDataTag ON RawData (tag);
			CREATE INDEX IF NOT EXISTS RawDataDs ON RawData (ds);
			CREATE INDEX IF NOT EXISTS RawDataTs ON RawData (ts);
			CREATE INDEX IF NOT EXISTS RawDataExpAt ON RawData (expat);
		COMMIT;
	`)
	if err != nil {
		log.Error("Error in DB: " + err.Error())
		panic(err)
	}
}

func (sq *SqliteTSS) CloseStorage() {
	sq.db.Close()
	sq.isRunning = false
}

func (sq *SqliteTSS) Save(dataSource string, data map[string]*proto.TSPoints, expirationMillis uint64) {
	sb := strings.Builder{}
	now := utils.GetNowMillis()
	expireAt := now + expirationMillis

	sb.WriteString("BEGIN TRANSACTION;")
	for tag, values := range data {
		for ts, value := range values.Points {
			sb.WriteString(fmt.Sprintf("INSERT INTO RawData VALUES(\"%s\", \"%s\", %d, %f, %d);", dataSource, tag, ts, value, expireAt))
		}
	}
	sb.WriteString("COMMIT;")

	rq := sb.String()
	_, err := sq.db.Exec(rq)
	if err != nil {
		log.Error("Error in DB in Save: " + err.Error())
	}
}

func (sq *SqliteTSS) Retrieve(dataSource string, tags []string, fromTimestamp uint64, toTimestamp uint64) map[string]*proto.TSPoints {
	ans := make(map[string]*proto.TSPoints)

	for _, tag := range tags {
		ansForTag := make(map[uint64]float64)
		rq := fmt.Sprintf("SELECT ts, value FROM RawData WHERE ds = \"%s\" AND tag = \"%s\" AND ts >= %d AND ts <= %d", dataSource, tag, fromTimestamp, toTimestamp)
		rows, err := sq.db.Query(rq)
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
	rq := fmt.Sprintf("SELECT min(ts), max(ts) FROM RawData WHERE ds = \"%s\";", dataSource)
	res, err := sq.db.Query(rq)
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
	rq := fmt.Sprintf("BEGIN TRANSACTION;DELETE FROM RawData WHERE expat <= %d;COMMIT;", now)
	_, err := sq.db.Exec(rq)

	if err != nil {
		log.Error("Error in DB in expirationCycle: " + err.Error())
	}
}
