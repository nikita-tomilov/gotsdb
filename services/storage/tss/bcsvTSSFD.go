package tss

import (
	"bufio"
	"encoding/binary"
	"github.com/btcsuite/btcutil/base58"
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/utils"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type BCSVTSforDatasource struct {
	DatasourcePath string
}

func (dsd *BCSVTSforDatasource) Init() {
	os.MkdirAll(dsd.DatasourcePath, os.ModePerm)
}

func (dsd *BCSVTSforDatasource) GetData(tags []string, fromTimestamp uint64, toTimestamp uint64) map[string]*pb.TSPoints {
	ans := make(map[string]*pb.TSPoints)
	for _, tag := range tags {
		data, _ := dsd.getDataFromFile(dsd.filenameForTag(tag), fromTimestamp, toTimestamp)
		ans[tag] = &pb.TSPoints{Points: data}
	}
	return ans
}

func (dsd *BCSVTSforDatasource) SaveData(data map[string]*pb.TSPoints, expiration uint64) {
	expireAt := utils.GetNowMillis() + expiration
	if expiration == 0 {
		expireAt = 0
	}
	for tag, values := range data {
		f, err := os.OpenFile(dsd.filenameForTag(tag), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		writer := bufio.NewWriter(f)
		for ts, value := range values.Points {
			if _, err := writer.Write(dsd.encodebCSVLine(ts, value, expireAt)); err != nil {
				log.Println(err)
			}
		}
		_ = writer.Flush()
		_ = f.Close()
	}
}

func (dsd *BCSVTSforDatasource) SaveDataBatch(data []*pb.TSPoint, expiration uint64) {
	ans := make(map[string]*pb.TSPoints)
	converted := make(map[string]map[uint64]float64)
	for _, point := range data {
		_, exists := converted[point.Tag]
		if !exists {
			converted[point.Tag] = make(map[uint64]float64)
		}
		converted[point.Tag][point.Timestamp] = point.Value
	}
	for tag, dataForTag := range converted {
		ans[tag] = &pb.TSPoints{Points:dataForTag}
	}
	dsd.SaveData(ans, expiration)
}

func (dsd *BCSVTSforDatasource) Availability(fromTimestamp uint64, toTimestamp uint64) []*pb.TSAvailabilityChunk {
	ansMin := uint64(math.MaxUint64)
	ansMax := uint64(0)

	for _, file := range dsd.getAllFiles() {
		data, _ := dsd.getDataFromFile(file, 0, uint64(math.MaxUint64))
		for ts, _ := range data {
			ansMin = utils.Min(ansMin, ts)
			ansMax = utils.Max(ansMax, ts)
		}
	}

	if ansMax < ansMin {
		ans := make([]*pb.TSAvailabilityChunk, 0)
		return ans
	}

	ansMin = utils.Max(fromTimestamp, ansMin)
	ansMax = utils.Min(toTimestamp, ansMax)

	ans := make([]*pb.TSAvailabilityChunk, 1)
	ans[0] = &pb.TSAvailabilityChunk{FromTimestamp: ansMin, ToTimestamp: ansMax}
	return ans
}

func (dsd *BCSVTSforDatasource) filenameForTag(tag string) string {
	return dsd.DatasourcePath + "/" + base58.Encode([]byte(tag)) + ".bcsv"
}

func (dsd *BCSVTSforDatasource) encodebCSVLine(ts uint64, value float64, expireAt uint64) []byte {
	entry := make([]byte, 24)
	binary.LittleEndian.PutUint64(entry, ts)
	binary.LittleEndian.PutUint64(entry[8:], math.Float64bits(value))
	binary.LittleEndian.PutUint64(entry[16:], expireAt)
	return entry
}

func (dsd *BCSVTSforDatasource) decodebCSVLine(entry []byte) (uint64, float64, uint64) {
	ts := binary.LittleEndian.Uint64(entry[0:8])
	val := utils.ByteToFloat64(entry[8:16])
	expAt := binary.LittleEndian.Uint64(entry[16:24])
	return ts, val, expAt
}

func (dsd *BCSVTSforDatasource) getDataFromFile(filePath string, from uint64, to uint64) (map[uint64]float64, map[uint64]uint64) {
	values := make(map[uint64]float64)
	expAts := make(map[uint64]uint64)
	now := utils.GetNowMillis()

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return values, expAts
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	buf := make([]byte, 24)
	reader := bufio.NewReader(file)

	_, err = io.ReadFull(reader, buf)
	for err == nil {
		ts, val, expAt := dsd.decodebCSVLine(buf)
		withinRange := (from <= ts) && (ts <= to)
		if withinRange && ((expAt > now) || (expAt == 0)) {
			values[ts] = val
			expAts[ts] = expAt
		}
		_, err = io.ReadFull(reader, buf)
	}

	return values, expAts
}

func (dsd *BCSVTSforDatasource) getAllFiles() []string {
	var ans []string
	err := filepath.Walk(dsd.DatasourcePath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".bcsv") {
			ans = append(ans, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return ans
}

func (dsd *BCSVTSforDatasource) ExpirationCycle() {
	for _, file := range dsd.getAllFiles() {
		data, expAts := dsd.getDataFromFile(file, 0, math.MaxUint64)
		timestamps := make([]uint64, 0, len(data))
		for k := range data {
			timestamps = append(timestamps, k)
		}
		sort.Slice(timestamps, func(i, j int) bool {
			return timestamps[i] < timestamps[j]
		})

		f, err := os.OpenFile(file + ".copy", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		writer := bufio.NewWriter(f)
		for _, ts := range timestamps {
			value := data[ts]
			expireAt := expAts[ts]
			if _, err := writer.Write(dsd.encodebCSVLine(ts, value, expireAt)); err != nil {
				log.Println(err)
			}
		}
		_ = writer.Flush()
		_ = f.Close()

		err = os.Rename(file + ".copy", file)
		if err != nil {
			panic(err)
		}
	}
}
