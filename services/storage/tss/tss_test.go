package tss

import (
	"fmt"
	pb "github.com/nikita-tomilov/gotsdb/proto"
	"github.com/nikita-tomilov/gotsdb/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTSS_BasicFunctionsWork1(t *testing.T) {
	//given
	storages := buildStorages()
	for _, s := range storages {
		func() {
			defer s.CloseStorage()
			dataToStore := buildData()
			//when
			s.Save(testDataSource, dataToStore, 15000)
			retrievedData := s.Retrieve(testDataSource, []string{testTag1, testTag2}, may040520, may050520)
			//then
			assert.Equal(t, dataToStore[testTag1], retrievedData[testTag1], "Data for tag1 for 04.05-05.05 should be same, storage %s", s.String())
			assert.Equal(t, make(map[uint64]float64), retrievedData[testTag2].Points, "Data for tag2 for 04.05-05.05 should be empty, storage %s", s.String())
		}()
	}
}

func TestTSS_BasicFunctionsWork2(t *testing.T) {
	//given
	storages := buildStorages()
	for _, s := range storages {
		func() {
			defer s.CloseStorage()
			dataToStore := buildData()
			//when
			s.Save(testDataSource, dataToStore, 1000)
			retrievedData := s.Retrieve(testDataSource, []string{testTag1, testTag2}, may050520, may060520)
			//then
			assert.Equal(t, make(map[uint64]float64), retrievedData[testTag1].Points, "Data for tag1 for 05.05-06.05 should be empty, storage %s", s.String())
			assert.Equal(t, dataToStore[testTag2], retrievedData[testTag2], "Data for tag2 for 05.05-06.05 should be same, storage %s", s.String())
		}()
	}
}

func TestTSS_BasicFunctionsWork3(t *testing.T) {
	//given
	storages := buildStorages()
	for _, s := range storages {
		func() {
			defer s.CloseStorage()
			dataToStore := buildData()
			//when
			s.Save(testDataSource, dataToStore, 1000)
			retrievedData := s.Retrieve(testDataSource, []string{testTag1, testTag2}, may040520, may060520)
			//then
			assert.Equal(t, dataToStore[testTag1], retrievedData[testTag1], "Data for tag1 for 04.05-06.05 should be same, storage %s", s.String())
			assert.Equal(t, dataToStore[testTag2], retrievedData[testTag2], "Data for tag2 for 04.05-06.05 should be empty, storage %s", s.String())
			//when
			avail := s.Availability(testDataSource, 0, may060520+10000)
			//then
			assert.Equal(t, 1, len(avail))
			assert.Equal(t, &pb.TSAvailabilityChunk{FromTimestamp: may040520, ToTimestamp: may050520 + 2000}, avail[0])
		}()
	}
}

func Test_ExpirationWorks1(t *testing.T) {
	//given
	storages := buildStorages()
	for _, s := range storages {
		func() {
			defer s.CloseStorage()
			dataToStore := buildData()
			//when
			s.Save(testDataSource, dataToStore, 1000)
			time.Sleep(time.Second * 3)
			retrievedData := s.Retrieve(testDataSource, []string{testTag1, testTag2}, may040520, may060520)
			//then
			assert.Equal(t, make(map[uint64]float64), retrievedData[testTag1].Points, "Data for tag1 should be expired, storage %s", s.String())
			assert.Equal(t, make(map[uint64]float64), retrievedData[testTag2].Points, "Data for tag2 should be expired, storage %s", s.String())
			//when
			avail := s.Availability(testDataSource, 0, may060520+10000)
			//then
			assert.Equal(t, 0, len(avail), "Availability is incorrect for %s", s.String())
		}()
	}
}

func Test_ExpirationWorks2(t *testing.T) {
	//given
	storages := buildStorages()
	for _, s := range storages {
		func() {
			defer s.CloseStorage()
			dataToStore := buildData()
			//when
			s.Save(testDataSource, dataToStore, 15000)
			time.Sleep(time.Second * 3)
			retrievedData := s.Retrieve(testDataSource, []string{testTag1, testTag2}, may040520, may060520)
			//then
			assert.Equal(t, dataToStore[testTag1], retrievedData[testTag1], "Data for tag1 should not be expired, storage %s", s.String())
			assert.Equal(t, dataToStore[testTag2], retrievedData[testTag2], "Data for tag2 should not be expired, storage %s", s.String())
			//when
			avail := s.Availability(testDataSource, 0, may060520+10000)
			//then
			assert.Equal(t, 1, len(avail), "Availability is incorrect for %s", s.String())
			assert.Equal(t, &pb.TSAvailabilityChunk{FromTimestamp: may040520, ToTimestamp: may050520 + 2000}, avail[0])
		}()
	}
}

func buildStorages() []TimeSeriesStorage {
	inMem := buildInMemStorage()
	qL := buildQlStorage()
	lsm := buildLSMStorage()
	return toArray(inMem, qL, lsm)
}

func toArray(items ...TimeSeriesStorage) []TimeSeriesStorage {
	return items
}

func buildInMemStorage() *InMemTSS {
	s := InMemTSS{periodBetweenWipes: time.Second * 1}
	s.InitStorage()
	return &s
}

func buildQlStorage() *QlBasedPersistentTSS {
	idx += 1
	s := QlBasedPersistentTSS{Path: fmt.Sprintf("/tmp/gotsdb_test/test%d%d", utils.GetNowMillis(), idx), periodBetweenWipes: time.Second * 1}
	s.InitStorage()
	return &s
}

func buildLSMStorage() *LSMTSS {
	idx += 1
	s := LSMTSS{Path: fmt.Sprintf("/tmp/gotsdb_test/test%d%d", utils.GetNowMillis(), idx), CommitlogFlushPeriodSeconds:1, MemtExpirationPeriodSeconds:1}
	s.InitStorage()
	return &s
}

func buildData() map[string]*pb.TSPoints {
	m := make(map[string]*pb.TSPoints)

	dataForTag1 := make(map[uint64]float64)
	dataForTag1[may040520] = 42.0
	dataForTag1[may040520+1000] = 69.0
	m[testTag1] = &pb.TSPoints{Points: dataForTag1}

	dataForTag2 := make(map[uint64]float64)
	dataForTag2[may050520+1000] = 42.0
	dataForTag2[may050520+2000] = 69.0
	m[testTag2] = &pb.TSPoints{Points: dataForTag2}

	return m
}

const testDataSource string = "test-ds"
const testTag1 string = "test-tag-1"
const testTag2 string = "test-tag-2"
const may040520 = 1588550400000
const may050520 = 1588636800000
const may060520 = 1588723200000
var idx = 0
