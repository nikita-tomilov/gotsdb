package tss

import (
	"github.com/nikita-tomilov/golsm/commitlog"
	"github.com/nikita-tomilov/golsm/memt"
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/utils"
	"time"
)

type TSforDatasource struct {
	memt               *memt.Manager
	PeriodBetweenWipes time.Duration
}

func (dataSourceData *TSforDatasource) Init() {
	dataSourceData.memt = &memt.Manager{PerformExpirationEvery: dataSourceData.PeriodBetweenWipes}
	dataSourceData.memt.InitStorage()
}

func (dataSourceData *TSforDatasource) GetData(tags []string, fromTimestamp uint64, toTimestamp uint64) map[string]*pb.TSPoints {
	ans := make(map[string]*pb.TSPoints)
	for _, tag := range tags {
		memtFt := dataSourceData.memt.MemTableForTag(tag)
		data := memtFt.Retrieve(fromTimestamp, toTimestamp)
		ans[tag] = convertEntriesToTSP(data)
	}
	return ans
}

func convertEntriesToTSP(measurements []memt.Entry) *pb.TSPoints {
	ans := make(map[uint64]float64)
	for _, m := range measurements {
		ans[m.Timestamp] = utils.ByteToFloat64(m.Value)
	}
	return &pb.TSPoints{Points: ans}
}

func (dataSourceData *TSforDatasource) SaveData(data map[string]*pb.TSPoints, expiration uint64) {
	expireAt := utils.GetNowMillis() + expiration
	if expiration == 0 {
		expireAt = 0
	}
	for tag, values := range data {
		memtFt := dataSourceData.memt.MemTableForTag(tag)

		convData := convertTSPtoEntries(values, tag, expireAt)
		memtFt.MergeWithCommitlog(convData)
	}
}

func convertTSPtoEntries(points *pb.TSPoints, tag string, expireAt uint64) []commitlog.Entry {
	ans := make([]commitlog.Entry, len(points.Points))
	i := 0
	for ts, val := range points.Points {
		ans[i] = commitlog.Entry{Key: []byte(tag), Timestamp: ts, Value: utils.Float64ToByte(val), ExpiresAt:expireAt}
		i++
	}
	return ans
}


func (dataSourceData *TSforDatasource) Availability(fromTimestamp uint64, toTimestamp uint64) []*pb.TSAvailabilityChunk {
	ansMin, ansMax := dataSourceData.memt.Availability()

	ansMin = utils.Max(fromTimestamp, ansMin)
	ansMax = utils.Min(toTimestamp, ansMax)

	if (ansMin == 0) && (ansMax == 0) {
		return []*pb.TSAvailabilityChunk{}
	}

	ans := make([]*pb.TSAvailabilityChunk, 1)
	ans[0] = &pb.TSAvailabilityChunk{FromTimestamp: ansMin, ToTimestamp: ansMax}
	return ans
}
