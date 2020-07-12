package tss

type Void struct{}

type TimeSeriesStorage interface {
	InitStorage()
	Save(dataSource string, data map[string]map[uint64]float64, expirationMillis uint64)
	Retrieve(dataSource string, tags []string, fromTimestamp uint64, toTimestamp uint64) map[string]map[uint64]float64
}
