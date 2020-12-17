package tss
/*
import (
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/nikita-tomilov/gotsdb/utils"
	"testing"
	"time"
)

func TestStorageConverter_CloneAlreadySavedFiles(t *testing.T) {
	log.LoadConfiguration("../../../config/log4go.json")
	engineFrom := LSMTSS{Path: fmt.Sprintf("/home/hotaro/go/src/github.com/nikita-tomilov/gotsdb/testdata/benchmark_read/lsm"),
		CommitlogFlushPeriodSeconds: 1, CommitlogMaxEntries: 10, MemtExpirationPeriodSeconds: 1, MemtMaxEntriesPerTag: 100, MemtPrefetchSeconds: 120}
	engineTo := QlBasedPersistentTSS{Path: "/home/hotaro/go/src/github.com/nikita-tomilov/gotsdb/testdata/benchmark_read/ql", periodBetweenWipes: time.Second * 10000}
	//engineTo := CSVTSS{Path: "/home/hotaro/go/src/github.com/nikita-tomilov/gotsdb/testdata/benchmark_read/csv", periodBetweenWipes: time.Second * 10000}
	//engineTo := InMemTSS{periodBetweenWipes: time.Second * 10000}

	engineFrom.InitStorage()
	engineTo.InitStorage()

	const dataSource = "whatever"
	CloneAlreadySavedFiles(&engineFrom, &engineTo, dataSource, engineFrom.GetTags(dataSource))

	availChunks := engineFrom.Availability(dataSource, 0, utils.GetNowMillis() * 2)
	avail := availChunks[0]
	tags := engineFrom.GetTags(dataSource)

	for _, tag := range tags {
		data := engineFrom.Retrieve(dataSource, []string{tag}, avail.FromTimestamp, avail.ToTimestamp)
		data2 := engineTo.Retrieve(dataSource, []string{tag}, avail.FromTimestamp, avail.ToTimestamp)
		dataPoints := len(data[tag].Points)
		dataPoints2 := len(data2[tag].Points)
		if dataPoints != dataPoints2 {
			fmt.Printf("tag '%s' mismatch: %d points in %s, %d points in %s\n", tag, dataPoints, engineFrom.String(), dataPoints2, engineTo.String())
		} else {
			fmt.Printf("tag '%s' copied successfully\n", tag)
		}
	}

	engineFrom.CloseStorage()
	engineTo.CloseStorage()
}*/