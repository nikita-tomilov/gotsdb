package tss

import (
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/nikita-tomilov/gotsdb/utils"
	"testing"
	"time"
)

func unixTsToString(ts uint64) string {
	tm := time.Unix(int64(ts / 1000), 0)
	return tm.String()
}

func TestLSMTSS_CanReadAlreasySavedFiles(t *testing.T) {
	//given
	log.LoadConfiguration("../../../config/log4go.json")
	s := LSMTSS{Path: fmt.Sprintf("/home/hotaro/go/src/github.com/nikita-tomilov/gotsdb/testdata"), CommitlogFlushPeriodSeconds: 1, CommitlogMaxEntries: 10, MemtExpirationPeriodSeconds: 1, MemtMaxEntriesPerTag: 100, MemtPrefetchSeconds: 120}
	const dataSource = "whatever"
	s.InitStorage()
	func() {
		defer s.CloseStorage()
		//when
		availChunks := s.Availability(dataSource, 0, utils.GetNowMillis())
		//then
		for _, chunk := range availChunks {
			fmt.Printf(" - %s to %s\n", unixTsToString(chunk.FromTimestamp), unixTsToString(chunk.ToTimestamp))
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

//func TestLSMTSS_CanCloneAlreasySavedFiles(t *testing.T) {
func LSMTSS_CanCloneAlreasySavedFiles(t *testing.T) {
	//given
	log.LoadConfiguration("../../../config/log4go.json")
	s := LSMTSS{Path: fmt.Sprintf("/home/hotaro/go/src/github.com/nikita-tomilov/gotsdb/testdata"), CommitlogFlushPeriodSeconds: 1, CommitlogMaxEntries: 10, MemtExpirationPeriodSeconds: 1, MemtMaxEntriesPerTag: 100, MemtPrefetchSeconds: 120}
	s2 := SqliteTSS{Path:"/home/hotaro/go/src/github.com/nikita-tomilov/gotsdb/testdata"}
	const dataSource = "whatever"
	s.InitStorage()
	s2.InitStorage()
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
		//when
		tags := s.GetTags(dataSource)
		//then
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
