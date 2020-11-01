package tss

import (
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/utils"
	"github.com/pkg/profile"
	"math"
	"math/rand"
	"testing"
	"time"
)

//for everything:
//GOGC=off go test 'github.com/nikita-tomilov/gotsdb/services/storage/tss' -test.run== -bench=. -benchtime=100x
//for only reading:

/*
GOGC=off GO111MODULE=off go test 'github.com/nikita-tomilov/gotsdb/services/storage/tss' -bench=BenchmarkDataReading -test.run== -benchtime=100x
*/

//https://stackoverflow.com/questions/16935965/how-to-run-test-cases-in-a-specified-file

func BenchmarkDataReading(b *testing.B) {
	storages := BuildStoragesForBenchmark()
	log.Close()

	dataFrom := utils.GetNowMillis()
	dataTo := dataFrom + 60*60*1000
	tagsCount := 10
	data := buildDummyDataForBenchmark(tagsCount, dataFrom, dataTo)
	ds := "whatever"
	for _, storage := range storages {

		println("=============\n\nSaving data to " + storage.String())
		saveStarted := time.Now()
		storage.Save(ds, data, 0)
		saveWas := time.Since(saveStarted)

		println("Saving data was for " + saveWas.String())
		time.Sleep(10 * time.Second)

		benchmarkName := "DataRead on " + storage.String()
		println("Starting " + benchmarkName)
		b.Run(benchmarkName, func(b *testing.B) {
			p := profile.Start(profile.CPUProfile, profile.ProfilePath("../../../"), profile.NoShutdownHook)
			for i := 0; i < b.N; i++ {
				from := randomTs(dataFrom+20, dataFrom+(dataTo-dataFrom)/2)
				to := randomTs(from, dataTo)
				if to-from <= 10 {
					from -= 10
				}
				d := storage.Retrieve(ds, []string{"tag1", "tag2", "tag3"}, from, to)
				if len(d) != 3 {
					panic("tags mismatch")
				}
			}
			p.Stop()
		})

		println("Finished for storage " + storage.String() + "\n\n")
		storage.CloseStorage()
	}
}

func BenchmarkDataReading_LSMvsSQLite(b *testing.B) {
	storages := BuildStoragesForBenchmarkLSMvsSQLite()
	log.Close()

	dataFrom := utils.GetNowMillis()
	millisInHour := 60 * 60 * 1000
	dataTo := dataFrom + uint64(millisInHour)
	tagsCount := 10
	data := buildDummyDataForBenchmark(tagsCount, dataFrom, dataTo)
	ds := "whatever"
	requestSizes := []time.Duration{time.Second * 10, time.Second * 60, time.Minute * 2, time.Minute * 5, time.Minute * 10, time.Minute * 30}

	for _, storage := range storages {

		saveToStorage(storage, ds, data)

		for _, requestSize := range requestSizes {

			avail := storage.Availability(ds, dataFrom, dataTo)
			fmt.Printf("Expected avail: %d to %d\nActual avail:\n", dataFrom, dataTo)
			for _, chunk := range avail {
				fmt.Printf("- %d to %d\n", chunk.FromTimestamp, chunk.ToTimestamp)
			}

			benchmarkName := "DataRead on " + storage.String() + " for " + requestSize.String()
			b.Run(benchmarkName, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					from, to := randomTimeRange(dataFrom, dataTo, uint64(requestSize.Milliseconds()))
					d := storage.Retrieve(ds, []string{"tag1", "tag2", "tag3"}, from, to)
					if len(d) != 3 {
						panic("tags mismatch")
					}
					//for tag, values := range d {
					//	if !withinDelta(float64(len(values.Points)), float64(requestSize.Milliseconds()/1000), 0.1) {
					//		print(fmt.Sprintf("count mismatch for %s: expected about %d, got %d\n", tag, requestSize.Milliseconds()/1000, len(values.Points)))
					//	}
					//}
				}
			})

		}
	}
}

func BenchmarkDataWriting(b *testing.B) {
	storages := BuildStoragesForBenchmark()
	log.Close()

	ds := "whatever"
	for _, storage := range storages {
		benchmarkName := "DataWrite on " + storage.String()
		println("Starting " + benchmarkName)
		b.Run(benchmarkName, func(b *testing.B) {
			p := profile.Start(profile.CPUProfile, profile.ProfilePath("../../../"), profile.NoShutdownHook)
			for i := 0; i < b.N; i++ {
				dataFrom := i * 1000
				dataTo := dataFrom + 999
				tagsCount := 100
				data := buildDummyDataForBenchmark(tagsCount, uint64(dataFrom), uint64(dataTo))
				storage.Save(ds, data, 0)
			}
			p.Stop()
		})

		println("Finished for storage " + storage.String() + "\n\n")
	}
}

func buildDummyDataForBenchmark(tagsCount int, tsFrom uint64, tsTo uint64) map[string]*proto.TSPoints {
	tags := make([]string, tagsCount)
	for i := 0; i < tagsCount; i++ {
		tags[i] = fmt.Sprintf("tag%d", i)
	}
	ans := make(map[string]*proto.TSPoints)
	for _, tag := range tags {
		data := make(map[uint64]float64)
		for ts := tsFrom; ts < tsTo; ts += 1000 {
			data[ts] = rand.Float64()
		}
		ans[tag] = &proto.TSPoints{Points: data}
	}
	return ans
}

func saveToStorage(storage TimeSeriesStorage, ds string, data map[string]*proto.TSPoints) {
	println("=============\n\nSaving data to " + storage.String())
	saveStarted := time.Now()
	storage.Save(ds, data, 0)
	saveWas := time.Since(saveStarted)

	println("Saving data was for " + saveWas.String())
	time.Sleep(15 * time.Second)
}

func randomTs(from uint64, to uint64) uint64 {
	return uint64(rand.Float64()*float64(to-from) + float64(from))
}

func randomTimeRange(from uint64, to uint64, width uint64) (uint64, uint64) {
	fromAns := randomTs(from, to-2*width)
	toAns := from + width
	return fromAns, toAns
}

func withinDelta(a float64, b float64, percentage float64) bool {
	epsilon := (a + b) / 2 * percentage
	return math.Abs(a-b) < epsilon
}
