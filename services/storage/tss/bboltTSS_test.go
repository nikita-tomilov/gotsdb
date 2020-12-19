package tss

import (
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/nikita-tomilov/gotsdb/utils"
	"testing"
	"time"
)

func TestBboltTSS_CanReadAlreasySavedFiles(t *testing.T) {
	//given
	log.LoadConfiguration("../../../config/log4go.json")
	s := BboltTSS{Path: "/home/hotaro/go/src/github.com/nikita-tomilov/gotsdb/testdata/benchmark_read/bbolt", periodBetweenWipes: time.Second * 10000}
	const dataSource = "whatever"
	s.InitStorage()
	func() {
		defer s.CloseStorage()
		//when
		availChunks := s.Availability(dataSource, 0, utils.GetNowMillis() * 2)
		//then
		for _, chunk := range availChunks {
			fmt.Printf(" - %s to %s\n", utils.UnixTsToString(chunk.FromTimestamp), utils.UnixTsToString(chunk.ToTimestamp))
		}
		avail := availChunks[0]
		//when
		tags := []string{"tag1", "tag2", "tag3"}
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