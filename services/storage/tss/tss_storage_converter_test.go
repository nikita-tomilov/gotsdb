package tss

/*
func TestStorageConverter_CloneAlreadySavedFiles(t *testing.T) {
	log.LoadConfiguration("../../../config/log4go.json")
	s := LSMTSS{Path: fmt.Sprintf("/home/hotaro/go/src/github.com/nikita-tomilov/gotsdb/testdata/benchmark_read/lsm"),
		CommitlogFlushPeriodSeconds: 1, CommitlogMaxEntries: 10, MemtExpirationPeriodSeconds: 1, MemtMaxEntriesPerTag: 100, MemtPrefetchSeconds: 120}
	//s2 := QlBasedPersistentTSS{Path: "/home/hotaro/go/src/github.com/nikita-tomilov/gotsdb/testdata/benchmark_read/ql", periodBetweenWipes: time.Second * 10000}
	s2 := CSVTSS{Path: "/home/hotaro/go/src/github.com/nikita-tomilov/gotsdb/testdata/benchmark_read/csv", periodBetweenWipes: time.Second * 10000}

	s.InitStorage()
	s2.InitStorage()

	const dataSource = "whatever"
	CloneAlreasySavedFiles(&s, &s2, dataSource, s.GetTags(dataSource))

	s.CloseStorage()
	s2.CloseStorage()
}
*/