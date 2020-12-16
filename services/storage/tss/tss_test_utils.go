package tss

import (
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/nikita-tomilov/gotsdb/utils"
	"time"
)

func BuildStoragesForTesting() []TimeSeriesStorage {
	inMem := buildInMemStorage()
	qL := buildQlStorage()
	lsm := buildLSMStorage()
	sQ := buildSqliteStorage()
	csv := buildCSVStorage()
	return toArray(inMem, qL, lsm, sQ, csv)
}

func BuildStoragesForBenchmark(path string) []TimeSeriesStorage {
	inmem := buildInMemStorageForBenchmark()
	lsm := buildLSMStorageForBenchmark(path + "/lsm")
	CloneAlreadySavedFiles(lsm, inmem, "whatever", lsm.GetTags("whatever"))
	sQ := buildSqliteStorageForBenchmark(path + "/sqlite")
	//ql := buildQlStorageForBenchmark(path + "/ql")
	csv := buildCSVStorageForBenchmark(path + "/csv")
	return toArray(inmem, csv, lsm, sQ)
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
	s := LSMTSS{Path: fmt.Sprintf("/tmp/gotsdb_test/test%d%d", utils.GetNowMillis(), idx), CommitlogFlushPeriodSeconds: 1, CommitlogMaxEntries: 10, MemtExpirationPeriodSeconds: 1, MemtMaxEntriesPerTag: 100, MemtPrefetchSeconds:120}
	s.InitStorage()
	return &s
}

func buildLSMStorageForBenchmark(path string) *LSMTSS {
	idx += 1
	s := LSMTSS{Path: path, CommitlogFlushPeriodSeconds: 5, CommitlogMaxEntries: 100, MemtExpirationPeriodSeconds: 3000, MemtMaxEntriesPerTag: 1000, MemtPrefetchSeconds:120}
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

func buildCSVStorageForBenchmark(path string) *CSVTSS {
	idx += 1
	s := CSVTSS{Path: path, periodBetweenWipes: time.Hour * 1024}
	s.InitStorage()
	return &s
}

var idx = 0

func CloneAlreadySavedFiles(s TimeSeriesStorage, s2 TimeSeriesStorage, dataSource string, tags []string) {
	//given
	func() {
		defer s.CloseStorage()
		defer s2.CloseStorage()
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
	}()
	log.Close()
}
