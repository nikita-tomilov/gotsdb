package dummytss

import (
	log "github.com/jeanphorn/log4go"
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"time"
)

//TODO: treemap?

type TSPointWithExpiration struct {
	value    float64
	expireAt uint64
}

type TSforTag struct {
	tag string
	data map[uint64]TSPointWithExpiration
}

func (tagData *TSforTag) Init(tag string) {
	tagData.tag = tag
	tagData.data = make(map[uint64]TSPointWithExpiration)
}

func (tagData *TSforTag) GetData(fromTimestamp uint64, toTimestamp uint64) *pb.TSPoints {
	ans := make(map[uint64]float64)
	for timestamp, point := range tagData.data {
		if (timestamp >= fromTimestamp) && (timestamp <= toTimestamp) {
			ans[timestamp] = point.value
		}
	}
	return &pb.TSPoints{Points: ans}
}

func (tagData *TSforTag) SaveData(data *pb.TSPoints, expirationMillis uint64) {
	now := getNowMillis()
	expireAt := now + expirationMillis
	for timestamp, value := range data.Points {
		tagData.data[timestamp] = TSPointWithExpiration{value, expireAt}
	}
}

func (tagData *TSforTag) ExpirationCycle() {
	now := getNowMillis()
	for ts, point := range tagData.data {
		if point.expireAt <= now {
			log.Debug("expiring point for %s ts %d as it expires at %d and now it is %d", tagData.tag, ts, point.expireAt, now)
			delete(tagData.data, ts)
		}
	}
}

func getNowMillis() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}
