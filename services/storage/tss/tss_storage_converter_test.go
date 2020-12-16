package tss

import (
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/nikita-tomilov/gotsdb/utils"
)

func CloneAlreadySavedFiles(s TimeSeriesStorage, s2 TimeSeriesStorage, dataSource string, tags []string) {
	//given
	func() {
		defer s.CloseStorage()
		defer s2.CloseStorage()
		//when
		availChunks := s.Availability(dataSource, 0, utils.GetNowMillis())
		//then
		for _, chunk := range availChunks {
			fmt.Printf(" - %s to %s\n", unixTsToString(chunk.FromTimestamp), unixTsToString(chunk.ToTimestamp))
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