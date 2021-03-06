package tss

import (
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/nikita-tomilov/gotsdb/utils"
	"testing"
)

func TestLSMTSS_CanReadAlreasySavedFiles(t *testing.T) {
	//given
	log.LoadConfiguration("../../../config/log4go.json")
	s := LSMTSS{Path: fmt.Sprintf("/home/hotaro/go/src/github.com/nikita-tomilov/gotsdb/testdata/benchmark_read/lsm"), CommitlogFlushPeriodSeconds: 1, CommitlogMaxEntries: 10, MemtExpirationPeriodSeconds: 1, MemtMaxEntriesPerTag: 100, MemtPrefetchSeconds: 120}
	const dataSource = "whatever"
	s.InitStorage()
	func() {
		defer s.CloseStorage()
		//when
		availChunks := s.Availability(dataSource, 0, utils.GetNowMillis())
		//then
		for _, chunk := range availChunks {
			fmt.Printf(" - %s to %s\n", utils.UnixTsToString(chunk.FromTimestamp), utils.UnixTsToString(chunk.ToTimestamp))
		}
		avail := availChunks[0]
		//when
		tags := s.GetTags(dataSource)
		//then
		for _, tag := range tags {
			//when
			data := s.Retrieve(dataSource, []string{tag}, avail.FromTimestamp, avail.ToTimestamp)
			//then
			fmt.Printf("tag '%s' has %d points\n", tag, len(data[tag].Points))
		}
	}()
	log.Close()
}