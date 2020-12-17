package tss

import (
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/utils"
	"math"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkDataReading(b *testing.B) {
	storages := BuildStoragesForBenchmark("/home/hotaro/go/src/github.com/nikita-tomilov/gotsdb/testdata/benchmark_read", true)
	log.Close()

	ds := "whatever"
	requestSizes := []time.Duration{
		time.Second * 5,
		time.Second * 10,
		time.Second * 15,
		time.Second * 20,
		time.Second * 25,
		time.Second * 30,
		time.Second * 45,
		time.Second * 60,
		time.Minute * 2,
		time.Minute * 3,
		time.Minute * 4,
		time.Minute * 5,
		time.Minute * 10,
		time.Minute * 15,
		time.Minute * 20,
		time.Minute * 25,
		time.Minute * 30,
		time.Minute * 45,
		time.Minute * 60,
		time.Minute * 75,
		time.Minute * 90,
		time.Minute * 105,
		time.Minute * 120,
		time.Minute * 135}

	for _, storage := range storages {
		avail := storage.Availability(ds, 0, 2*utils.GetNowMillis())
		dataFrom := avail[0].FromTimestamp
		dataTo := avail[0].ToTimestamp

		for _, requestSize := range requestSizes {
			benchmarkName := fmt.Sprintf("DataRead on %s for %s |%d|", storage.String(), requestSize.String(), int(requestSize.Seconds()))
			b.Run(benchmarkName, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					from, to := randomTimeRange(dataFrom, dataTo, uint64(requestSize.Milliseconds()))
					d := storage.Retrieve(ds, []string{"tag1", "tag2", "tag3"}, from, to)
					if len(d) != 3 {
						panic("tags mismatch")
					}
					for tag, dataForTag := range d {
						expected := int(requestSize.Milliseconds() / 1000)
						actual := len(dataForTag.Points)
						if !withinDelta(expected, actual, 0.1) {
							fmt.Printf("mismatch on tag %s; expected %d got %d on timerange %s - %s\n", tag, expected, actual, utils.UnixTsToString(from), utils.UnixTsToString(to))
						}
					}
				}
			})

		}

		storage.CloseStorage()
	}
}

func BenchmarkLatestDataReading(b *testing.B) {
	storages := BuildStoragesForBenchmark("/home/hotaro/go/src/github.com/nikita-tomilov/gotsdb/testdata/benchmark_read", true)
	log.Close()

	ds := "whatever"
	requestSizes := []time.Duration{
		time.Second * 5,
		time.Second * 10,
		time.Second * 15,
		time.Second * 20,
		time.Second * 25,
		time.Second * 30,
		time.Second * 45,
		time.Second * 60,
		time.Second * 75,
		time.Second * 90}
	for _, storage := range storages {
		avail := storage.Availability(ds, 0, 2*utils.GetNowMillis())
		dataTo := avail[0].ToTimestamp
		dataFrom := dataTo - 3*60*1000
		for _, requestSize := range requestSizes {
			benchmarkName := fmt.Sprintf("LatestDataRead on %s for %s |%d|", storage.String(), requestSize.String(), int(requestSize.Seconds()))
			b.Run(benchmarkName, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					from, to := randomTimeRange(dataFrom, dataTo, uint64(requestSize.Milliseconds()))
					d := storage.Retrieve(ds, []string{"tag1", "tag2", "tag3"}, from, to)
					if len(d) != 3 {
						panic("tags mismatch")
					}
					for tag, dataForTag := range d {
						expected := int(requestSize.Milliseconds() / 1000)
						actual := len(dataForTag.Points)
						if !withinDelta(expected, actual, 0.1) {
							fmt.Printf("mismatch on tag %s; expected %d got %d on timerange %s - %s\n", tag, expected, actual, utils.UnixTsToString(from), utils.UnixTsToString(to))
						}
					}
				}
			})

		}

		storage.CloseStorage()
	}
}

func BenchmarkLinearDataWriting(b *testing.B) {
	storages := BuildStoragesForBenchmark(fmt.Sprintf("/tmp/gotsdb/benchdata%d", utils.GetNowMillis()), false)
	log.Close()

	ds := "whatever"
	requestSizes := []time.Duration{
		time.Second * 5,
		time.Second * 10,
		time.Second * 15,
		time.Second * 20,
		time.Second * 25,
		time.Second * 30,
		time.Second * 45,
		time.Second * 60,
		time.Minute * 2,
		time.Minute * 3,
		time.Minute * 4,
		time.Minute * 5,
		time.Minute * 10,
		time.Minute * 15,
		time.Minute * 20,
		time.Minute * 25,
		time.Minute * 30,
	}
	for _, storage := range storages {
		var timeFrom = uint64(0)
		for _, requestSize := range requestSizes {
			benchmarkName := fmt.Sprintf("LinearDataWrite on %s for %s |%d|", storage.String(), requestSize.String(), int(requestSize.Seconds()))
			b.Run(benchmarkName, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					timeTo := timeFrom + uint64(requestSize.Milliseconds())
					randomData := buildDummyDataForBenchmark(10, timeFrom, timeTo)
					storage.Save(ds, randomData, 0)
					timeFrom += 2 * uint64(requestSize.Milliseconds())
				}
			})
		}
		storage.CloseStorage()
	}
}

func BenchmarkRandomDataWriting(b *testing.B) {
	storages := BuildStoragesForBenchmark(fmt.Sprintf("/tmp/gotsdb/benchdata%d", utils.GetNowMillis()), false)
	log.Close()

	ds := "whatever"
	requestSizes := []time.Duration{
		time.Second * 5,
		time.Second * 10,
		time.Second * 15,
		time.Second * 20,
		time.Second * 25,
		time.Second * 30,
		time.Second * 45,
		time.Second * 60,
		time.Minute * 2,
		time.Minute * 3,
		time.Minute * 4,
		time.Minute * 5,
		time.Minute * 10,
		time.Minute * 15,
		time.Minute * 20,
		time.Minute * 25,
		time.Minute * 30,
	}
	for _, storage := range storages {
		now := utils.GetNowMillis()
		for _, requestSize := range requestSizes {
			benchmarkName := fmt.Sprintf("RandomDataWrite on %s for %s |%d|", storage.String(), requestSize.String(), int(requestSize.Seconds()))
			b.Run(benchmarkName, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					timeFrom := randomTs(0, now)
					timeTo := timeFrom + uint64(requestSize.Milliseconds())
					randomData := buildDummyDataForBenchmark(10, timeFrom, timeTo)
					storage.Save(ds, randomData, 0)
				}
			})
		}
		storage.CloseStorage()
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

func randomTs(from uint64, to uint64) uint64 {
	return uint64(rand.Float64()*float64(to-from) + float64(from))
}

func randomTimeRange(from uint64, to uint64, width uint64) (uint64, uint64) {
	fromAns := randomTs(from, to-width)
	toAns := fromAns + width
	return fromAns, toAns
}

func withinDelta(a int, b int, percentage float64) bool {
	epsilon := float64(a+b) * percentage
	return math.Abs(float64(a)-float64(b)) <= epsilon
}
