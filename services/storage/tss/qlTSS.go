package tss

import (
	"errors"
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/utils"
	"math"
	"modernc.org/ql"
	"os"
	"strings"
	"time"
)

type QlBasedPersistentTSS struct {
	Path               string `summer.property:"tss.filePath|/tmp/gotsdb/tss"`
	dbFilePath         string
	periodBetweenWipes time.Duration
	commonImpl         AbstractSQLTSS
}

func (qp *QlBasedPersistentTSS) InitStorage() {
	_ = os.MkdirAll(qp.Path, os.ModePerm)

	if qp.periodBetweenWipes == 0 * time.Second {
		qp.periodBetweenWipes = time.Second * 5
	}

	qp.dbFilePath = qp.Path + "/db.bin"
	db, err := ql.OpenFile(qp.dbFilePath, &ql.Options{CanCreate: true, FileFormat: 2, RemoveEmptyWAL: true})
	if err != nil {
		panic("Unable to instantiate db " + err.Error())
	}
	ctx := &ql.TCtx{}

	sqlWrapper := QlWrapperImpl{ctx: ctx, db: db}
	qp.commonImpl = AbstractSQLTSS{sqlWrapper: &sqlWrapper, periodBetweenWipes: qp.periodBetweenWipes}
	qp.commonImpl.Init()
}

func (qp *QlBasedPersistentTSS) CloseStorage() {
	qp.commonImpl.Close()
}

func (qp *QlBasedPersistentTSS) Save(dataSource string, data map[string]*proto.TSPoints, expirationMillis uint64) {
	qp.commonImpl.Save(dataSource, data, expirationMillis)
}

func (qp *QlBasedPersistentTSS) Retrieve(dataSource string, tags []string, fromTimestamp uint64, toTimestamp uint64) map[string]*proto.TSPoints {
	return qp.commonImpl.Retrieve(dataSource, tags, fromTimestamp, toTimestamp)
}

func (qp *QlBasedPersistentTSS) Availability(dataSource string, fromTimestamp uint64, toTimestamp uint64) []*proto.TSAvailabilityChunk {
	return qp.commonImpl.Availability(dataSource, fromTimestamp, toTimestamp)
}

func (qp *QlBasedPersistentTSS) String() string {
	return fmt.Sprintf("QlBasedPersistentTSS on file %s", qp.dbFilePath)
}

// Ql specific part

/*
~/go/src/github.com/nikita-tomilov/gotsdb/testdata/benchmark_read/ql$ ql -db db.bin -fld 'select min(ts), max(ts) from RawData;'
*/

type QlWrapperImpl struct {
	ctx *ql.TCtx
	db  *ql.DB
}

func (q *QlWrapperImpl) InitDatabase() {
	q.Execute(`
		BEGIN TRANSACTION;
			CREATE TABLE IF NOT EXISTS measurement_meta (id int, tag string);
			CREATE TABLE IF NOT EXISTS measurement_data_sources (id int, data_source string);
			CREATE TABLE IF NOT EXISTS measurements (data_source_key int, meta_key int, ts uint64, value float64, expire_at uint64);
			CREATE INDEX IF NOT EXISTS idx_measurements_meta_key ON measurements(meta_key);
			CREATE INDEX IF NOT EXISTS idx_measurements_ds_key ON measurements(data_source_key);
			CREATE INDEX IF NOT EXISTS idx_measurements_ts ON measurements(ts);
			CREATE INDEX IF NOT EXISTS idx_measurements_expire_at ON measurements(expire_at);
		COMMIT;
	`)
}

func (q *QlWrapperImpl) Execute(query string) {
	_, _, err := q.db.Run(q.ctx, query)
	if err != nil {
		log.Error("Error in DB: " + err.Error())
		panic(err)
	}
}

func (q *QlWrapperImpl) CreateMetaKey(tag string) {
	sb := strings.Builder{}
	sb.WriteString("BEGIN TRANSACTION;")
	sb.WriteString(fmt.Sprintf("INSERT INTO measurement_meta (tag) VALUES(\"%s\");", tag))
	sb.WriteString("COMMIT;")
	rq := sb.String()
	q.Execute(rq)
}

func (q *QlWrapperImpl) GetMetaKey(query string) (uint, error) {
	rq := fmt.Sprintf("SELECT id() FROM measurement_meta WHERE tag = \"%s\"", query)
	res, _, err := q.db.Run(q.ctx, rq)
	if err != nil {
		log.Error("Error in DB: " + err.Error())
		panic(err)
	}
	key := uint(0)

	for _, ress := range res {
		rows, _ := ress.Rows(math.MaxInt64, 0)
		if rows == nil {
			return key, errors.New("not found")
		}
		for _, row := range rows {
			if row[0] != nil {
				key = uint(row[0].(int64))
			}
		}
	}

	return key, nil
}

func (q *QlWrapperImpl) CreateDataSourceKey(dataSource string) {
	sb := strings.Builder{}
	sb.WriteString("BEGIN TRANSACTION;")
	sb.WriteString(fmt.Sprintf("INSERT INTO measurement_data_sources (data_source) VALUES(\"%s\");", dataSource))
	sb.WriteString("COMMIT;")
	rq := sb.String()
	q.Execute(rq)
}

func (q *QlWrapperImpl) GetDataSourceKey(query string) (uint, error) {
	rq := fmt.Sprintf("SELECT id() FROM measurement_data_sources WHERE data_source = \"%s\"", query)
	res, _, err := q.db.Run(q.ctx, rq)
	if err != nil {
		log.Error("Error in DB: " + err.Error())
		panic(err)
	}
	key := uint(0)

	for _, ress := range res {
		rows, _ := ress.Rows(math.MaxInt64, 0)
		if rows == nil {
			return key, errors.New("not found")
		}
		for _, row := range rows {
			if row[0] != nil {
				key = uint(row[0].(int64))
			}
		}
	}

	return key, nil
}

func (q *QlWrapperImpl) GetTwoTimestamps(query string) (uint64, uint64) {
	res, _, err := q.db.Run(q.ctx, query)
	min := uint64(math.MaxUint64)
	max := uint64(0)
	if err != nil {
		log.Error("Error in DB in Availability: " + err.Error())
		return min, max
	}
	for _, ress := range res {
		rows, _ := ress.Rows(math.MaxInt64, 0)
		for _, row := range rows {
			if (row[0] != nil) && (row[1] != nil) {
				min = utils.Min(min, row[0].(uint64))
				max = utils.Max(max, row[1].(uint64))
			}
		}
	}
	return min, max
}

func (q *QlWrapperImpl) GetMeasurementsForTag(query string) map[uint64]float64 {
	ansForTag := make(map[uint64]float64)
	res, _, err := q.db.Run(q.ctx, query)
	now := utils.GetNowMillis()
	if err != nil {
		log.Error("Error in DB in Retrieve: " + err.Error())
	} else {
		for _, ress := range res {
			rows, _ := ress.Rows(math.MaxInt64, 0)
			for _, row := range rows {
				ts := row[0].(uint64)
				val := row[1].(float64)
				expAt := row[2].(uint64)
				if (expAt == 0) || (expAt > now) {
					ansForTag[ts] = val
				}
			}
		}
	}
	return ansForTag
}

func (q *QlWrapperImpl) DeleteOnExpiration() {
	now := utils.GetNowMillis()
	rq := fmt.Sprintf("BEGIN TRANSACTION;DELETE FROM measurements WHERE (expire_at != 0) AND (expire_at < %d);COMMIT;", now)
	_, _, err := q.db.Run(q.ctx, rq)

	if err != nil {
		log.Error("Error in DB in expirationCycle: " + err.Error())
	}
}

func (q *QlWrapperImpl) Close() {
	_ = q.db.Close()
}
