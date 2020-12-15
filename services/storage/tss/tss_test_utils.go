package tss

import (
	"fmt"
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

func BuildStoragesForBenchmark() []TimeSeriesStorage {
	idx += 1
	inMem := buildInMemStorage()
	qL := buildQlStorage()
	lsm := buildLSMStorageForBenchmark(fmt.Sprintf("/tmp/gotsdb_test/test%d%d", utils.GetNowMillis(), idx))
	sQ := buildSqliteStorage()
	return toArray(inMem, qL, lsm, sQ)
}


func BuildStoragesForBenchmarkLSMvsSQLite(path string) []TimeSeriesStorage {
	lsm := buildLSMStorageForBenchmark(path)
	sQ := buildSqliteStorageForBenchmark(path)
	return toArray(lsm, sQ)
}

func toArray(items ...TimeSeriesStorage) []TimeSeriesStorage {
	return items
}

func buildInMemStorage() *InMemTSS {
	s := InMemTSS{periodBetweenWipes: time.Second * 1}
	s.InitStorage()
	return &s
}

func buildQlStorage() *QlBasedPersistentTSS {
	idx += 1
	s := QlBasedPersistentTSS{Path: fmt.Sprintf("/tmp/gotsdb_test/test%d%d", utils.GetNowMillis(), idx), periodBetweenWipes: time.Second * 1}
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
	s := CSVTSS{Path: fmt.Sprintf("/tmp/gotsdb_test/csvtest%d%d", utils.GetNowMillis(), idx), periodBetweenWipes: time.Hour * 1024}
	s.InitStorage()
	return &s
}

var idx = 0
