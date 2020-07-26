package tss

import (
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/nikita-tomilov/gotsdb/proto"
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/utils"
	"math"
	"modernc.org/ql"
	"os"
	"strings"
	"time"
)

type QlBasedPersistentTSS struct {
	ctx  *ql.TCtx
	db   *ql.DB
	Path string `summer.property:"ts.filePath|/tmp/gotsdb/tss"`
	dbFilePath string
	periodBetweenWipes time.Duration
	isRunning bool
}

func (qp *QlBasedPersistentTSS) InitStorage() {
	_ = os.MkdirAll(qp.Path, os.ModePerm)
	qp.dbFilePath = qp.Path+"/db.bin"
	db, err := ql.OpenFile(qp.dbFilePath, &ql.Options{CanCreate: true, FileFormat: 1})
	if err != nil {
		panic("Unable to instantiate db " + err.Error())
	}
	qp.ctx = &ql.TCtx{}
	qp.db = db
	qp.createTableIfNotExists()

	qp.isRunning = true
	if qp.periodBetweenWipes == 0 * time.Second {
		qp.periodBetweenWipes = time.Second * 5
	}
	go func(qp *QlBasedPersistentTSS) {
		for qp.isRunning {
			qp.expirationCycle()
			time.Sleep(qp.periodBetweenWipes)
		}
	}(qp)
}

func (qp *QlBasedPersistentTSS) createTableIfNotExists() {
	_, _, err := qp.db.Run(qp.ctx, `
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

func (qp *QlBasedPersistentTSS) CloseStorage() {
	qp.db.Close()
	qp.isRunning = false
}

func (qp *QlBasedPersistentTSS) Save(dataSource string, data map[string]*proto.TSPoints, expirationMillis uint64) {
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

	rq := sb.String();
	_, _, err := qp.db.Run(qp.ctx, rq)
	if err != nil {
		log.Error("Error in DB in Save: " + err.Error())
	}
}

func (qp *QlBasedPersistentTSS) Retrieve(dataSource string, tags []string, fromTimestamp uint64, toTimestamp uint64) map[string]*proto.TSPoints {
	ans := make(map[string]*proto.TSPoints)

	for _, tag := range tags {
		ansForTag := make(map[uint64]float64)
		rq := fmt.Sprintf("SELECT * FROM RawData WHERE ds = \"%s\" AND tag = \"%s\" AND ts >= %d AND ts <= %d", dataSource, tag, fromTimestamp, toTimestamp)
		res, _, err := qp.db.Run(qp.ctx, rq)
		if err != nil {
			log.Error("Error in DB in Retrieve: " + err.Error())
		} else {
			for _, ress := range res {
				rows, _ := ress.Rows(math.MaxInt64, 0)
				for _, row := range rows {
					ts := row[2].(uint64)
					val := row[3].(float64)
					ansForTag[ts] = val
				}
			}
		}
		ans[tag] = &pb.TSPoints{Points: ansForTag}
	}

	return ans
}

func (qp *QlBasedPersistentTSS) Availability(dataSource string, fromTimestamp uint64, toTimestamp uint64) []*proto.TSAvailabilityChunk {
	rq := fmt.Sprintf("SELECT min(ts), max(ts) FROM RawData WHERE ds = \"%s\";", dataSource)
	res, _, err := qp.db.Run(qp.ctx, rq)
	min := uint64(math.MaxUint64)
	max := uint64(0)
	if err != nil {
		log.Error("Error in DB in Availability: " + err.Error())
		ans := make([]*proto.TSAvailabilityChunk, 0)
		return ans
	} else {
		for _, ress := range res {
			rows, _ := ress.Rows(math.MaxInt64, 0)
			for _, row := range rows {
				if (row[0] != nil) && (row[1] != nil) {
					min = utils.Min(min, row[0].(uint64))
					max = utils.Max(max, row[1].(uint64))
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
	ans[0] = &proto.TSAvailabilityChunk{FromTimestamp:min, ToTimestamp:max}
	return ans
}

func (qp *QlBasedPersistentTSS) String() string {
	return fmt.Sprintf("QlBasedPersistentTSS on file %s", qp.dbFilePath)
}

func (qp *QlBasedPersistentTSS) expirationCycle() {
	now := utils.GetNowMillis()
	rq := fmt.Sprintf("BEGIN TRANSACTION;DELETE FROM RawData WHERE expat <= %d;COMMIT;", now)
	_, _, err := qp.db.Run(qp.ctx, rq)

	if err != nil {
		log.Error("Error in DB in expirationCycle: " + err.Error())
	}
}