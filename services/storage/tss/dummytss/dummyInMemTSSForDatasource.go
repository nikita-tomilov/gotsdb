package dummytss

import pb "github.com/nikita-tomilov/gotsdb/proto"

type TSforDatasource struct {
	data map[string]TSforTag
}

func (dataSourceData *TSforDatasource) Init() {
	dataSourceData.data = make(map[string]TSforTag)
}

func (dataSourceData *TSforDatasource) GetData(tags []string, fromTimestamp uint64, toTimestamp uint64) map[string]*pb.TSPoints {
	ans := make(map[string]*pb.TSPoints)
	for _, tag := range tags {
		if dataSourceData.contains(tag) {
			ans[tag] = dataSourceData.dataForTag(tag).GetData(fromTimestamp, toTimestamp)
		} else {
			ans[tag] = &pb.TSPoints{Points: make(map[uint64]float64)} //TODO: throw exception?
		}
	}
	return ans
}

func (dataSourceData *TSforDatasource) SaveData(data map[string]*pb.TSPoints, expiration uint64) {
	for tag, values := range data {
		if !dataSourceData.contains(tag) {
			dataForTag := TSforTag{}
			dataForTag.Init(tag)
			dataSourceData.data[tag] = dataForTag
		}
		dataSourceData.dataForTag(tag).SaveData(values, expiration)
	}
}

func (dataSourceData *TSforDatasource) ExpirationCycle() {
	for _, data := range dataSourceData.data {
		data.ExpirationCycle()
	}
}

func (dataSourceData *TSforDatasource) contains(tag string) bool {
	_, found := dataSourceData.data[tag]
	return found
}

func (dataSourceData *TSforDatasource) dataForTag(tag string) *TSforTag {
	dataForTag := dataSourceData.data[tag]
	return &dataForTag
}