package dummytss

import "time"

//TODO: treemap?

type TSPointWithExpiration struct {
	value    float64
	expireAt uint64
}

type TSforTag struct {
	data map[uint64]TSPointWithExpiration
}

func (tagData *TSforTag) Init() {
	tagData.data = make(map[uint64]TSPointWithExpiration)
}

func (tagData *TSforTag) GetData(fromTimestamp uint64, toTimestamp uint64) map[uint64]float64 {
	ans := make(map[uint64]float64)
	for timestamp, point := range tagData.data {
		if (timestamp >= fromTimestamp) && (timestamp <= toTimestamp) {
			ans[timestamp] = point.value
		}
	}
	return ans
}

func (tagData *TSforTag) SaveData(data map[uint64]float64, expirationMillis uint64) {
	now := getNowMillis()
	expireAt := now + expirationMillis
	for timestamp, value := range data {
		tagData.data[timestamp] = TSPointWithExpiration{value, expireAt}
	}
}

func (tagData *TSforTag) ExpirationCycle() {
	now := getNowMillis()
	for ts, point := range tagData.data {
		if point.expireAt <= now {
			delete(tagData.data, ts)
		}
	}
}

func getNowMillis() uint64 {
	return uint64(time.Now().UnixNano() / 1000)
}
