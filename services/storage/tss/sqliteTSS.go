package tss

import (
	"database/sql"
	"errors"
	"fmt"
	log "github.com/jeanphorn/log4go"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"math"
	"os"
	"strings"
	"time"
)

type SqliteTSS struct {
	Path               string `summer.property:"tss.filePath|/tmp/gotsdb/tss"`
	dbFilePath         string
	periodBetweenWipes time.Duration
	commonImpl         AbstractSQLTSS
}

func (sq *SqliteTSS) InitStorage() {
	_ = os.MkdirAll(sq.Path, os.ModePerm)

	if sq.periodBetweenWipes == 0*time.Second {
		sq.periodBetweenWipes = time.Second * 5
	}

	sq.dbFilePath = sq.Path + "/db.bin"
	db, err := gorm.Open(sqlite.Open(sq.dbFilePath), &gorm.Config{})
	if err != nil {
		panic("Unable to instantiate db " + err.Error())
	}

	sqlWrapper := SqliteWrapperImpl{db: db}
	sq.commonImpl = AbstractSQLTSS{sqlWrapper: &sqlWrapper, periodBetweenWipes: sq.periodBetweenWipes}
	sq.commonImpl.Init()
}

func (sq *SqliteTSS) CloseStorage() {
	sq.commonImpl.Close()
}

func (sq *SqliteTSS) Save(dataSource string, data map[string]*proto.TSPoints, expirationMillis uint64) {
	sq.commonImpl.Save(dataSource, data, expirationMillis)
}

func (sq *SqliteTSS) SaveBatch(dataSource string, data []*proto.TSPoint, expirationMillis uint64) {
	sq.commonImpl.SaveBatch(dataSource, data, expirationMillis)
}

func (sq *SqliteTSS) Retrieve(dataSource string, tags []string, fromTimestamp uint64, toTimestamp uint64) map[string]*proto.TSPoints {
	return sq.commonImpl.Retrieve(dataSource, tags, fromTimestamp, toTimestamp)
}

func (sq *SqliteTSS) Availability(dataSource string, fromTimestamp uint64, toTimestamp uint64) []*proto.TSAvailabilityChunk {
	return sq.commonImpl.Availability(dataSource, fromTimestamp, toTimestamp)
}

func (sq *SqliteTSS) String() string {
	return fmt.Sprintf("SqliteTSS on file %s", sq.dbFilePath)
}

type SqliteWrapperImpl struct {
	db *gorm.DB
}

func (s *SqliteWrapperImpl) InitDatabase() {
	err := s.db.AutoMigrate(&Measurement{})
	if err != nil {
		log.Error("Error in DB: " + err.Error())
		panic(err)
	}
	err = s.db.AutoMigrate(&MeasurementMeta{})
	if err != nil {
		log.Error("Error in DB: " + err.Error())
		panic(err)
	}
	err = s.db.AutoMigrate(&MeasurementDataSource{})
	if err != nil {
		log.Error("Error in DB: " + err.Error())
		panic(err)
	}
}

func (s *SqliteWrapperImpl) Execute(query string) {
	db, _ := s.db.DB()
	_, err := db.Exec(query)
	if err != nil {
		log.Error("Error in DB in insertKeyByDsAndTag: " + err.Error())
	}
}

func (s *SqliteWrapperImpl) DeleteOnExpiration() {
	now := utils.GetNowMillis()
	err := s.db.Begin().Delete(Measurement{}, "expire_at > 0 AND expire_at <= ?", now).Commit().Error

	if err != nil {
		log.Error("Error in DB in expirationCycle: " + err.Error())
	}
}

func (s *SqliteWrapperImpl) CreateMetaKey(tag string) {
	sb := strings.Builder{}
	sb.WriteString("BEGIN TRANSACTION;")
	sb.WriteString(fmt.Sprintf("INSERT INTO measurement_meta (tag) VALUES(\"%s\");", tag))
	sb.WriteString("COMMIT;")
	rq := sb.String()
	db, _ := s.db.DB()
	_, err := db.Exec(rq)
	if err != nil {
		log.Error("Error in DB in insertKeyByDsAndTag: " + err.Error())
	}
}

func (s *SqliteWrapperImpl) GetMetaKey(tag string) (uint, error) {
	var meta MeasurementMeta
	_ = s.db.Where("tag = ?", tag).Find(&meta).Error
	if meta.Tag != tag {
		return 0, errors.New("not found")
	}
	return meta.Id, nil
}

func (s *SqliteWrapperImpl) CreateDataSourceKey(dataSource string) {
	sb := strings.Builder{}
	sb.WriteString("BEGIN TRANSACTION;")
	sb.WriteString(fmt.Sprintf("INSERT INTO measurement_data_sources (data_source) VALUES(\"%s\");", dataSource))
	sb.WriteString("COMMIT;")
	rq := sb.String()
	db, _ := s.db.DB()
	_, err := db.Exec(rq)
	if err != nil {
		log.Error("Error in DB in insertKeyByDsAndTag: " + err.Error())
	}
}

func (s *SqliteWrapperImpl) GetDataSourceKey(dataSource string) (uint, error) {
	var ds MeasurementDataSource
	_ = s.db.Where("data_source = ?", dataSource).Find(&ds).Error
	if ds.DataSource != dataSource {
		return 0, errors.New("not found")
	}
	return ds.Id, nil
}

func (s *SqliteWrapperImpl) GetTwoTimestamps(query string) (uint64, uint64) {
	res, err := s.db.Raw(query).Rows()
	min := uint64(math.MaxUint64)
	max := uint64(0)
	if err != nil {
		log.Error("Error in DB in Availability: " + err.Error())
		return min, max
	}
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
	return min, max
}

func (s *SqliteWrapperImpl) GetMeasurementsForTag(query string) map[uint64]float64 {
	ansForTag := make(map[uint64]float64)
	db, _ := s.db.DB()
	rows, err := db.Query(query)
	now := utils.GetNowMillis()
	if err != nil {
		log.Error("Error in DB in Retrieve: " + err.Error())
	} else {
		ts := uint64(0)
		val := 0.0
		expAt := uint64(0)
		for rows.Next() {
			err := rows.Scan(&ts, &val, &expAt)
			if err != nil {
				log.Error(err)
			} else {
				if (expAt == 0) || (expAt > now) {
					ansForTag[ts] = val
				}
			}
		}
	}
	return ansForTag
}

func (s *SqliteWrapperImpl) Close() {
	db, _ := s.db.DB()
	db.Close()
}
