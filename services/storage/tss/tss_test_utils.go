package tss

import (
	"fmt"
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/utils"
	"os"
	"time"
)

func BuildStoragesForTesting() []TimeSeriesStorage {
	inMem := buildInMemStorage()
	qL := buildQlStorage()
	lsm := buildLSMStorage()
	sQ := buildSqliteStorage()
	csv := buildCSVStorage()
	bcsv := buildBCSVStorage()
	bbolt := buildBboltStorage()
	return toArray(inMem, qL, lsm, sQ, csv, bcsv, bbolt)
}

func BuildStoragesForBenchmark(path string, readBenchmark bool) []TimeSeriesStorage {
	if !readBenchmark {
		os.MkdirAll(path + "/sqlite", os.ModePerm)
		os.MkdirAll(path + "/csv", os.ModePerm)
		os.MkdirAll(path + "/ql", os.ModePerm)
		os.MkdirAll(path + "/lsm", os.ModePerm)
		os.MkdirAll(path + "/bbolt", os.ModePerm)
	}
	inmem := buildInMemStorageForBenchmark()
	lsm := buildLSMStorageForBenchmark(path + "/lsm")
	if readBenchmark {
		CloneAlreadySavedFiles(lsm, inmem, "whatever", lsm.GetTags("whatever"))
	}
	sQ := buildSqliteStorageForBenchmark(path + "/sqlite")
	csv := buildCSVStorageForBenchmark(path + "/csv")
	bcsv := buildBCSVStorageForBenchmark(path + "/bcsv")
	qL := buildQlStorageForBenchmark(path + "/ql")
	bbolt := buildBboltStorageForBenchmark(path + "/bbolt")
	return toArray(inmem, csv, bcsv, lsm, sQ, qL, bbolt)
}

func toArray(items ...TimeSeriesStorage) []TimeSeriesStorage {
	return items
}

func buildInMemStorage() *InMemTSS {
	s := InMemTSS{periodBetweenWipes: time.Second * 1}
	s.InitStorage()
	return &s
}

func buildInMemStorageForBenchmark() *InMemTSS {
	s := InMemTSS{periodBetweenWipes: time.Second * 1024}
	s.InitStorage()
	return &s
}

func buildQlStorage() *QlBasedPersistentTSS {
	idx += 1
	s := QlBasedPersistentTSS{Path: fmt.Sprintf("/tmp/gotsdb_test/test%d%d", utils.GetNowMillis(), idx), periodBetweenWipes: time.Second * 1}
	s.InitStorage()
	return &s
}

func buildQlStorageForBenchmark(path string) *QlBasedPersistentTSS {
	idx += 1
	s := QlBasedPersistentTSS{Path: path, periodBetweenWipes: time.Second * 1024}
	s.InitStorage()
	return &s
}

func buildLSMStorage() *LSMTSS {
	idx += 1
	s := LSMTSS{Path: fmt.Sprintf("/tmp/gotsdb_test/test%d%d", utils.GetNowMillis(), idx), CommitlogFlushPeriodSeconds: 1, CommitlogMaxEntries: 10, MemtExpirationPeriodSeconds: 1, MemtMaxEntriesPerTag: 100, MemtPrefetchSeconds: 120}
	s.InitStorage()
	return &s
}

func buildLSMStorageForBenchmark(path string) *LSMTSS {
	idx += 1
	s := LSMTSS{Path: path, CommitlogFlushPeriodSeconds: 5, CommitlogMaxEntries: 100, MemtExpirationPeriodSeconds: 3000, MemtMaxEntriesPerTag: 1000, MemtPrefetchSeconds: 120}
	s.InitStorage()
	return &s
}

func buildSqliteStorage() *SqliteTSS {
	idx += 1
	s := SqliteTSS{Path: fmt.Sprintf("/tmp/gotsdb_test/test%d%d", utils.GetNowMillis(), idx), periodBetweenWipes: time.Second * 1}
	s.InitStorage()
	return &s
}

func buildSqliteStorageForBenchmark(path string) *SqliteTSS {
	idx += 1
	s := SqliteTSS{Path: path, periodBetweenWipes: time.Hour * 1024}
	s.InitStorage()
	return &s
}

func buildCSVStorage() *CSVTSS {
	idx += 1
	s := CSVTSS{Path: fmt.Sprintf("/tmp/gotsdb_test/csvtest%d%d", utils.GetNowMillis(), idx), periodBetweenWipes: time.Second * 1}
	s.InitStorage()
	return &s
}

func buildBCSVStorage() *BCSVTSS {
	idx += 1
	s := BCSVTSS{Path: fmt.Sprintf("/tmp/gotsdb_test/csvtest%d%d", utils.GetNowMillis(), idx), periodBetweenWipes: time.Second * 1}
	s.InitStorage()
	return &s
}

func buildCSVStorageForBenchmark(path string) *CSVTSS {
	idx += 1
	s := CSVTSS{Path: path, periodBetweenWipes: time.Hour * 1024}
	s.InitStorage()
	return &s
}

func buildBCSVStorageForBenchmark(path string) *BCSVTSS {
	idx += 1
	s := BCSVTSS{Path: path, periodBetweenWipes: time.Hour * 1024}
	s.InitStorage()
	return &s
}

func buildBboltStorage() *BboltTSS {
	idx += 1
	s := BboltTSS{Path: fmt.Sprintf("/tmp/gotsdb_test/bbolttest%d%d", utils.GetNowMillis(), idx), periodBetweenWipes: time.Second * 1}
	s.InitStorage()
	return &s
}

func buildBboltStorageForBenchmark(path string) *BboltTSS {
	idx += 1
	s := BboltTSS{Path: path, periodBetweenWipes: time.Hour * 1024}
	s.InitStorage()
	return &s
}

var idx = 0

func CloneAlreadySavedFiles(s TimeSeriesStorage, s2 TimeSeriesStorage, dataSource string, tags []string) {
	//when
	availChunks := s.Availability(dataSource, 0, utils.GetNowMillis())
	//then
	for _, chunk := range availChunks {
		fmt.Printf(" - %s to %s\n", utils.UnixTsToString(chunk.FromTimestamp), utils.UnixTsToString(chunk.ToTimestamp))
	}
	avail := availChunks[0]
	//when /then
	for _, tag := range tags {
		//when
		data := s.Retrieve(dataSource, []string{tag}, avail.FromTimestamp, avail.ToTimestamp)
		//then
		fmt.Printf("tag '%s' has %d points\n", tag, len(data[tag].Points))
		s2.Save(dataSource, data, 0)
	}
}


func ConvertToBatch(d map[string]*pb.TSPoints) []*pb.TSPoint {
	m := make([]*pb.TSPoint, 0)
	for tag, values := range d {
		for ts, val := range values.Points {
			m = append(m, &pb.TSPoint{Tag:tag, Timestamp:ts, Value:val})
		}
	}
	return m
}